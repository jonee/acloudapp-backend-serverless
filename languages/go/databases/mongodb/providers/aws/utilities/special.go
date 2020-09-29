// Created by Jonee Ryan Ty
// Copyright ACloudApp

/**
 * Special Utility functions
 */

package utilities

import (
	"context"
	"errors"
	"log"
	"os"

	"strings"
	"time"

	"encoding/json"
	"net/http"

	acaGoConfiguration "acloudapp.org/configuration"
	acaGoMongoDBModels "acloudapp.org/databases/mongodb/models"
	acaGoMongoDBAWSConfiguration "acloudapp.org/databases/mongodb/providers/aws/configuration"
	acaGoUtilities "acloudapp.org/utilities"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"

	jwt "github.com/dgrijalva/jwt-go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

// stage, bearer token, post and get parameters
func DoInit(mapStore map[string]interface{}, request events.APIGatewayProxyRequest, requireParameters bool) (bool, map[string]interface{}) {
	var err error
	hasError := false
	response := make(map[string]interface{})

	log.Println("FunctionName, FunctionVersion:", os.Getenv("AWS_LAMBDA_FUNCTION_NAME"), os.Getenv("AWS_LAMBDA_FUNCTION_VERSION"))

	// log.Println(request)

	// stage - dev or prod
	// stage := request.RequestContext.Stage
	stage := "dev"
	if request.Path != "" { // Resource or path would work
		stage = request.Path[1:] // remove first "/"
		stage = stage[0:strings.Index(stage, "/")]
	}
	mapStore["stage"] = stage
	log.Println("stage:", stage)

	// authorization bearer token from headers if any
	authHeader, _ := request.Headers["Authorization"]
	tokenString := ""
	if authHeader != "" {
		tokenString = acaGoUtilities.GetBearerToken(authHeader)
	}
	mapStore["tokenString"] = tokenString
	// log.Println("tokenString:", tokenString)

	mapStore["userAgent"], _ = request.Headers["User-Agent"]

	// post parameters
	var t map[string]interface{}
	if request.Body != "" {
		err = json.Unmarshal([]byte(request.Body), &t)
		if requireParameters && err != nil {
			hasError = true
			response = map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_PARAMETERS_REQUIRED", nil),
			}
		}
	} else {
		if requireParameters {
			hasError = true
			response = map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_PARAMETERS_REQUIRED", nil),
			}
		}
	}

	mapStore["t"] = t

	return hasError, response

	// DoInit
}

func DoCheckJwt(mapStore map[string]interface{}) (bool, map[string]interface{}) {
	var err error
	hasError := false
	response := make(map[string]interface{})

	tokenString := mapStore["tokenString"].(string)

	// check tokenString
	if tokenString == "" {
		return true, map[string]interface{}{
			"statusCode": http.StatusBadRequest,
			"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_JWT_REQUIRED", nil),
		}
	}

	// decode jwt
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("ERROR_JWT_INVALID_SIGNING_METHOD" + ": " + token.Header["alg"].(string))
		}

		return []byte(acaGoConfiguration.JWT_SECRET_STRING), nil
	})

	var userIdString string
	var userId primitive.ObjectID
	var loginLogIdString string
	var loginLogId primitive.ObjectID

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// log.Println(token)

		if token.Header["alg"].(string) != "HS256" {
			return true, map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_JWT_INVALID_SIGNING_METHOD", map[string]string{"alg": token.Header["alg"].(string)}),
			}
		}

		userIdString = claims["user_id"].(string)
		loginLogIdString = claims["login_log_id"].(string)
		exp := int64(claims["exp"].(float64))

		// check expiry
		if time.Now().Unix() > exp {
			return true, map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_JWT_EXPIRED", nil),
			}
		}

		// check bson ids
		userId, err = primitive.ObjectIDFromHex(userIdString)
		if err != nil {
			return true, map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_BAD_REQUEST", nil),
			}
		}

		loginLogId, err = primitive.ObjectIDFromHex(loginLogIdString)
		if err != nil {
			return true, map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_BAD_REQUEST", nil),
			}
		}

	} else {
		log.Println("ERROR", err)

		return true, map[string]interface{}{
			"statusCode": http.StatusBadRequest,
			"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_JWT_INVALID", nil),
		}
	}

	// should be done with simple validation still need to check db

	userCol := mapStore["userCol"].(*mongo.Collection)
	loginLogCol := mapStore["loginLogCol"].(*mongo.Collection)

	ctx := context.Background()

	f1 := bson.M{"_id": userId}
	var userObj acaGoMongoDBModels.User
	err = userCol.FindOne(ctx, f1).Decode(&userObj)
	if err != nil { // user not found
		return true, map[string]interface{}{
			"statusCode": http.StatusNotFound,
			"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_USER_NOT_FOUND", nil),
		}
	}
	mapStore["userObj"] = userObj

	f2 := bson.M{"_id": loginLogId}
	var loginLogObj acaGoMongoDBModels.LoginLog
	err = loginLogCol.FindOne(ctx, f2).Decode(&loginLogObj)
	if err != nil { // login log not found
		return true, map[string]interface{}{
			"statusCode": http.StatusNotFound,
			"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_LOGIN_LOG_NOT_FOUND", nil),
		}
	}
	mapStore["loginLogObj"] = loginLogObj

	// check if isblocked
	if userObj.IsBlocked {
		return true, map[string]interface{}{
			"statusCode": http.StatusBadRequest,
			"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_USER_IS_BLOCKED", nil),
		}
	}

	return hasError, response

	// DoCheckJwt
}

func DoEmail(mapStore map[string]interface{}, userObj acaGoMongoDBModels.User, emailType string, customParameterHolder string) {
	var err error

	if userObj.Email == "" {
		return
	}

	stage := mapStore["stage"].(string)

	language := "en"
	if userObj.Language != "" {
		language = userObj.Language
	}

	emailParams := map[string]string{"BRANDING": acaGoConfiguration.BRANDING, "username": userObj.Username}

	emailText := ""
	emailHTML := ""

	emailSubject := ""

	filter := bson.M{"language": language, "name": emailType}

	ctx := context.Background()
	emailTemplateCol := mapStore["emailTemplateCol"].(*mongo.Collection)

	var tmpEmailTemplateObj acaGoMongoDBModels.EmailTemplate
	err = emailTemplateCol.FindOne(ctx, filter).Decode(&tmpEmailTemplateObj)
	if err == nil { // object found
		emailSubject = tmpEmailTemplateObj.Subject
		emailText = tmpEmailTemplateObj.TemplateText
		emailHTML = tmpEmailTemplateObj.TemplateHTML
	}

	if emailType == "VALIDATION" {
		userLanguage := acaGoConfiguration.DEFAULT_LANGUAGE
		if userObj.Language != "" {
			userLanguage = userObj.Language
		}

		validationLink := strings.Replace(acaGoMongoDBAWSConfiguration.ACCOUNT_EMAIL_VALIDATION_LINK_PREFIX, "%%stage%%", stage, -1) + acaGoUtilities.EncodeB64(userObj.Id.Hex()) + "_" + userObj.ValidationSecret + "?language=" + userLanguage

		emailParams["validation_link"] = validationLink

	} else if emailType == "EMAIL_CHANGED_VALIDATION" {
		userLanguage := acaGoConfiguration.DEFAULT_LANGUAGE
		if userObj.Language != "" {
			userLanguage = userObj.Language
		}

		validationLink := strings.Replace(acaGoMongoDBAWSConfiguration.ACCOUNT_EMAIL_VALIDATION_LINK_PREFIX, "%%stage%%", stage, -1) + acaGoUtilities.EncodeB64(userObj.Id.Hex()) + "_" + userObj.ValidationSecret + "?language=" + userLanguage

		emailParams["validation_link"] = validationLink

	} else if emailType == "RESEND_VALIDATION" {
		userLanguage := acaGoConfiguration.DEFAULT_LANGUAGE
		if userObj.Language != "" {
			userLanguage = userObj.Language
		}

		validationLink := strings.Replace(acaGoMongoDBAWSConfiguration.ACCOUNT_EMAIL_VALIDATION_LINK_PREFIX, "%%stage%%", stage, -1) + acaGoUtilities.EncodeB64(userObj.Id.Hex()) + "_" + userObj.ValidationSecret + "?language=" + userLanguage

		emailParams["validation_link"] = validationLink

	} else if emailType == "FORGOT_PASSWORD" {
		emailParams["new_password"] = customParameterHolder

	} else if emailType == "MESSAGE_NOTIFICATION" {
		emailParams["from"] = customParameterHolder

	} else {
		log.Println("EMAIL ERROR UNKNOWN EMAIL TYPE")
	}

	if stage != "prod" {
		emailSubject = emailSubject + " (" + stage + ")"
	}

	// emailParams["subject"] = emailSubject

	emailText, emailHTML, emailSubject = acaGoUtilities.ProcessTemplate(emailParams, emailText, emailHTML, emailSubject)

	/*
		from := acaGoConfiguration.EMAIL_SMTP_FROM_ACCOUNT
		from_password := acaGoConfiguration.EMAIL_SMTP_FROM_PASSWORD
		to := userObj.Email

		body := emailMessage
		smtpServer := acaGoConfiguration.EMAIL_SMTP_SERVER
		smtpPort := acaGoConfiguration.EMAIL_SMTP_PORT
		err = SendEmailSmtp(from, from_password, to, emailSubject, body, smtpServer, smtpPort)
		if err != nil {
			log.Println("EMAIL ERROR", err)
		}
	*/

	Sender := acaGoConfiguration.EMAIL_FROM_ACCOUNT
	Recipient := userObj.Email
	Subject := emailSubject

	HtmlBody := emailHTML
	TextBody := emailText

	CharSet := "UTF-8"

	awsCreds := credentials.NewStaticCredentials(acaGoConfiguration.AWS_ACCESS_KEY_ID, acaGoConfiguration.AWS_SECRET_ACCESS_KEY, "")
	_, err = awsCreds.Get()
	if err != nil {
		log.Println("bad aws credentials")
	}

	awsCfg := aws.NewConfig().WithRegion(acaGoConfiguration.AWS_REGION).WithCredentials(awsCreds)
	awsSvc := ses.New(session.New(), awsCfg)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := awsSvc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
	} else {
		log.Println("Email Sent to address: " + Recipient)
		log.Println(result)
	}

	// DoEmail
}

func DoResponse(response map[string]interface{}, ct string) events.APIGatewayProxyResponse {
	AWSresponse := events.APIGatewayProxyResponse{
		StatusCode: response["statusCode"].(int),
		Body:       response["body"].(string),
	}

	if ct == "" {
		ct = "application/json"
	}

	AWSresponse.Headers = map[string]string{"Content-Type": ct}

	return AWSresponse

	// DoResponse
}
