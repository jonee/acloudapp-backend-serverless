// Created by Jonee Ryan Ty
// Copyright ACloudApp

// login

package main

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"net/http"

	acaGoConfiguration "acloudapp.org/configuration"
	acaGoMongoDBModels "acloudapp.org/databases/mongodb/models"
	acaGoMongoDBOpenWhiskUtilities "acloudapp.org/databases/mongodb/providers/openwhisk/utilities"
	acaGoMongoDBUtilities "acloudapp.org/databases/mongodb/utilities"
	acaGoUtilities "acloudapp.org/utilities"

	jwt "github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

var t0 = time.Now()
var mongoClient *mongo.Client

func main() {
	// AclOpen

	t1 := time.Now()
	log.Println("0-1 lambda life", (t1.UnixNano()-t0.UnixNano())/int64(time.Millisecond))

	mapStore := make(map[string]interface{})
	mapStore["last_milestone_time"] = t1

	var err error
	hasError := false

	hasError, response := acaGoMongoDBOpenWhiskUtilities.DoInit(mapStore, false)
	if hasError {
		acaGoMongoDBOpenWhiskUtilities.DoResponse(response, "")
		return
	}

	acaGoUtilities.PrintMilestone(mapStore, "do init")

	stage := mapStore["stage"].(string)
	t := mapStore["t"].(map[string]interface{})

	// parameters
	version, _ := t["version"].(string)
	applicationType, _ := t["application_type"].(string) // ios, android, react-native-ios, react-native-android, web-js, xamarin-wm10, ionic-tizen etc
	securityHash, _ := t["security_hash"].(string)       // security_hash - a formula of version, app_type, app_secret
	log.Println("version, applicationType, securityHash:", version, applicationType, acaGoUtilities.MaskString(securityHash, "*"))

	username, _ := t["username"].(string)
	email, _ := t["email"].(string)
	password, _ := t["password"].(string)
	// log.Println("username, email, password:", username, email, utilities.MaskString(password, "*"))

	// some dev cheat
	longExpiryFlag := ""
	if acaGoConfiguration.DEV_CHEAT_SLUG != "" {
		longExpiryFlag, _ = t[acaGoConfiguration.DEV_CHEAT_SLUG].(string)
	}

	// validate security hash TODO

	// validate parameters
	formErrors := make(map[string]interface{})
	hasError = false

	// validate username
	username = strings.TrimLeft(strings.TrimSpace(username), "@") // remove trailing @
	usernameErrors := make([](map[string]interface{}), 0)

	if username != "" {
		if len(username) < acaGoConfiguration.USERNAME_MIN_LENGTH { // min length
			usernameErrors = append(usernameErrors, map[string]interface{}{"message_key": "ERROR_USERNAME_MIN_LENGTH", "message_parameters": map[string]int{"USERNAME_MIN_LENGTH": acaGoConfiguration.USERNAME_MIN_LENGTH}})
			hasError = true
		} else {
			// regex
			usernameRe := regexp.MustCompile(acaGoConfiguration.REGEX_USERNAME_SERVER)
			if !usernameRe.MatchString(username) {
				usernameErrors = append(usernameErrors, map[string]interface{}{"message_key": "ERROR_USERNAME_REGEX"})
				hasError = true
			}
		}
	}

	if len(usernameErrors) > 0 {
		formErrors["username"] = usernameErrors
	}

	// validate email
	email = strings.TrimSpace(email)
	email = strings.Replace(email, " ", "+", -1) // since openwhisk changes + to space for parameters, we assume that a space in the middle is a +

	emailErrors := make([]map[string]interface{}, 0)

	if email != "" {
		// regex
		emailRe := regexp.MustCompile(acaGoConfiguration.REGEX_EMAIL)
		if !emailRe.MatchString(email) {
			emailErrors = append(emailErrors, map[string]interface{}{"message_key": "ERROR_EMAIL_REGEX"})
			hasError = true
		}
	}

	if len(emailErrors) > 0 {
		formErrors["email"] = emailErrors
	}

	// validate password
	passwordErrors := make([]map[string]interface{}, 0)

	if password == "" {
		passwordErrors = append(passwordErrors, map[string]interface{}{"message_key": "ERROR_PASSWORD_REQUIRED"})
		hasError = true
	} else if len(password) < acaGoConfiguration.PASSWORD_MIN_LENGTH {
		passwordErrors = append(passwordErrors, map[string]interface{}{"message_key": "ERROR_PASSWORD_MIN_LENGTH", "message_parameters": map[string]int{"PASSWORD_MIN_LENGTH": acaGoConfiguration.PASSWORD_MIN_LENGTH}})
		hasError = true
	}

	if len(passwordErrors) > 0 {
		formErrors["password"] = passwordErrors
	}

	if username == "" && email == "" {
		usernameErrors = append(usernameErrors, map[string]interface{}{"message_key": "ERROR_USERNAME_OR_EMAIL_REQUIRED"})
		formErrors["username"] = usernameErrors
		hasError = true
	}

	if hasError {
		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetFormErrorsJsonStringBodyReturn(false, "ERROR_VALIDATION", formErrors),
			},
			"",
		)
		return
	}

	acaGoUtilities.PrintMilestone(mapStore, "parameters")

	// should be done with simple validation still need to check db

	acaGoMongoDBUtilities.DoDBConnect(mapStore)

	acaGoUtilities.PrintMilestone(mapStore, "db")

	mongoClient = mapStore["mongoClient"].(*mongo.Client) // retrieve reference

	userCol := mapStore["userCol"].(*mongo.Collection)
	// loginLogCol := mapStore["loginLogCol"].(*mongo.Collection)

	var userObj acaGoMongoDBModels.User
	// foundUserObj := false

	// find username or email if already existing
	var filter bson.M
	if username != "" && email != "" {
		filter = bson.M{"$or": []bson.M{bson.M{"username": username}, bson.M{"email": email}}}
	} else if username != "" {
		filter = bson.M{"username": username}
	} else if email != "" {
		filter = bson.M{"email": email}
	}

	findOptions := options.Find()
	ctx := context.TODO()

	var results []acaGoMongoDBModels.User
	cur, err := userCol.Find(ctx, filter, findOptions)
	cur.All(ctx, &results)

	if err == nil { // object(s) found
		if len(results) > 1 {
			acaGoMongoDBOpenWhiskUtilities.DoResponse(
				map[string]interface{}{
					"statusCode": http.StatusBadRequest,
					"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_MULTIPLE_RESULTS", nil),
				},
				"",
			)
			return

		} else if len(results) == 0 {
			log.Println("hmm possible error 1")

			acaGoMongoDBOpenWhiskUtilities.DoResponse(
				map[string]interface{}{
					"statusCode": http.StatusNotFound,
					"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_USER_NOT_FOUND", nil),
				},
				"",
			)
			return

		} else { // len == 1
			userObj = results[0]
			// foundUserObj = true
		}

	} else {
		log.Println("ERROR", err)

		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusNotFound,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_USER_NOT_FOUND", nil),
			},
			"",
		)
		return
	}

	// check if validated
	if !userObj.IsEmailValidated {
		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_USER_SHOULD_VALIDATE", nil),
			},
			"",
		)
		return
	}

	// check if isblocked
	if userObj.IsBlocked {
		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_USER_IS_BLOCKED", nil),
			},
			"",
		)
		return
	}

	acaGoUtilities.PrintMilestone(mapStore, "dbs search")

	// check if password or temporary password matched
	// compare it with regular password hash
	err = bcrypt.CompareHashAndPassword([]byte(userObj.PasswordHash), []byte(userObj.PasswordSalt+password))
	log.Println("result 1", err)

	if err == nil {
		// success = true

	} else {
		// compare to passwordtemporaryhash
		err = bcrypt.CompareHashAndPassword([]byte(userObj.PasswordTemporaryHash), []byte(userObj.PasswordSalt+password))
		log.Println("result 2", err)
		if err == nil {
			if userObj.PasswordTemporaryExpiry.After(time.Now()) {
				// success = true
				userObj.PasswordHash = userObj.PasswordTemporaryHash
				userObj.PasswordTemporaryHash = ""
				// userObj.PasswordTemporaryExpiry = time.Now() // shall we unset? no so we have timestamp of (last) password request

				// userObj will be saved below

			} else {
				acaGoMongoDBOpenWhiskUtilities.DoResponse(
					map[string]interface{}{
						"statusCode": http.StatusBadRequest,
						"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_TEMPORARY_PASSWORD_EXPIRED", nil),
					},
					"",
				)
				return
			}

		} else {
			acaGoMongoDBOpenWhiskUtilities.DoResponse(
				map[string]interface{}{
					"statusCode": http.StatusBadRequest,
					"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_PASSWORD_IS_INCORRECT", nil),
				},
				"",
			)
			return
		}

	}

	acaGoUtilities.PrintMilestone(mapStore, "bcrypt compare")

	// all validations should now be done finally do save
	// create login log
	var loginLogObj acaGoMongoDBModels.LoginLog
	loginLogObj.UserId = userObj.Id
	loginLogObj.Username = userObj.Username
	loginLogObj.Email = userObj.Email
	loginLogObj.Version = version
	loginLogObj.ApplicationType = applicationType

	loginLogObj.Expiry = time.Now().Add(time.Hour * acaGoConfiguration.LOGGED_IN_EXPIRY)

	if (stage != "prod") && (longExpiryFlag == "1") {
		loginLogObj.Expiry = time.Now().Add(time.Hour * acaGoConfiguration.LOGGED_IN_EXPIRY_DEV_CHEAT)
	}

	loginLogObj.IsValid = true
	loginLogObj.Secret = acaGoUtilities.GetRandomString(8, "")
	loginLogObj.CreatedAt = time.Now()

	_, err = loginLogObj.Save(mapStore)
	if err != nil {
		log.Println("ERROR", err)

		log.Println("likely database down")

		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusInternalServerError,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_INTERNAL", nil),
			},
			"",
		)
		return
	}

	userObj.LoginCount++
	userObj.UpdatedAt = time.Now()

	_, err = userObj.Save(mapStore)
	if err != nil {
		log.Println("ERROR", err)

		log.Println("likely database down 2")

		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusInternalServerError,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_INTERNAL", nil),
			},
			"",
		)
		return
	}

	acaGoUtilities.PrintMilestone(mapStore, "db saving / create loginLog")

	// create jwt
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      userObj.Id.Hex(),
		"login_log_id": loginLogObj.Id.Hex(),
		"exp":          loginLogObj.Expiry.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	jwtTokenString, err := jwtToken.SignedString([]byte(acaGoConfiguration.JWT_SECRET_STRING))

	if err != nil {
		log.Println("ERROR", err)

		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusInternalServerError,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_INTERNAL", nil),
			},
			"",
		)
		return
	}

	ret := loginLogObj.ExportArrayPublic()
	ret["jwt_token"] = jwtTokenString
	ret["login_count"] = userObj.LoginCount

	/*
		if userObj.ProfilePictureLink != "" {
			ret["profile_picture_link"] = userObj.ProfilePictureLink
		}
	*/

	acaGoUtilities.PrintMilestone(mapStore, "jwt")

	log.Println("function exec", (time.Now().UnixNano()-t1.UnixNano())/int64(time.Millisecond))

	// return success
	acaGoMongoDBOpenWhiskUtilities.DoResponse(
		map[string]interface{}{
			"statusCode": http.StatusOK,
			"body":       acaGoUtilities.GetJsonStringBodyReturn(true, ret),
		},
		"",
	)
	return

	// main
}
