// Created by Jonee Ryan Ty
// Copyright ACloudApp

// email_validation

package main

import (
	"context"
	"log"
	"strings"
	"time"

	"net/http"

	acaGoConfiguration "acloudapp.org/configuration"
	acaGoMongoDBModels "acloudapp.org/databases/mongodb/models"
	acaGoMongoDBOpenWhiskUtilities "acloudapp.org/databases/mongodb/providers/openwhisk/utilities"
	acaGoMongoDBUtilities "acloudapp.org/databases/mongodb/utilities"
	acaGoUtilities "acloudapp.org/utilities"

	"github.com/avct/uasurfer"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
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

	// var err error
	hasError := false

	hasError, response := acaGoMongoDBOpenWhiskUtilities.DoInit(mapStore, false)
	if hasError {
		acaGoMongoDBOpenWhiskUtilities.DoResponse(response, "text/html")
		return
	}

	acaGoUtilities.PrintMilestone(mapStore, "do init")

	// t := mapStore["t"].(map[string]interface{})
	event := mapStore["event"].(map[string]interface{})

	// get parameters
	language, _ := event["language"].(string) // check for any in get parameters
	if language == "" {
		language = acaGoConfiguration.DEFAULT_LANGUAGE
	}
	log.Println("language:", language)

	// request path parameter
	oWPath := event["__ow_path"].(string) // /dev/account/email_validation/NWY2Yjg1ODZjNDFlYWIyMmFlN2YxOTRk_HTjeF1XH
	oWPathExplode := strings.Split(oWPath, "/")
	path := oWPathExplode[len(oWPathExplode)-1]
	// log.Println("path:", path)

	split := strings.Split(path, "_")
	if len(split) != 2 {
		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_BAD_REQUEST", nil, language),
			},
			"text/html",
		)
		return
	}

	mongoIdString := acaGoUtilities.DecodeB64(split[0])
	validationSecret := split[1]

	mongoId, err := primitive.ObjectIDFromHex(mongoIdString)
	// check bson ids
	if err != nil {
		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_BAD_REQUEST", nil, language),
			},
			"text/html",
		)
		return
	}

	acaGoUtilities.PrintMilestone(mapStore, "parameters")

	// should be done with simple parameter validation still need to check db

	acaGoMongoDBUtilities.DoDBConnect(mapStore)

	acaGoUtilities.PrintMilestone(mapStore, "db")

	mongoClient = mapStore["mongoClient"].(*mongo.Client) // retrieve reference

	userCol := mapStore["userCol"].(*mongo.Collection)
	// loginLogCol := mapStore["loginLogCol"].(*mongo.Collection)

	ctx := context.Background()

	// find user and validate parameters
	f := bson.M{"_id": mongoId}
	var userObj acaGoMongoDBModels.User
	err = userCol.FindOne(ctx, f).Decode(&userObj)
	if err == nil { // object found
		if userObj.ValidationSecret != validationSecret {
			acaGoMongoDBOpenWhiskUtilities.DoResponse(
				map[string]interface{}{
					"statusCode": http.StatusBadRequest,
					"body":       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_BAD_REQUEST", nil, language),
				},
				"text/html",
			)
			return

		} else if userObj.IsEmailValidated {
			acaGoMongoDBOpenWhiskUtilities.DoResponse(
				map[string]interface{}{
					"statusCode": http.StatusBadRequest,
					"body":       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_USER_ALREADY_VALIDATED", nil, language),
				},
				"text/html",
			)
			return
		}

	} else {
		// user not found
		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusNotFound,
				"body":       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_USER_NOT_FOUND", nil, language),
			},
			"text/html",
		)
		return
	}

	acaGoUtilities.PrintMilestone(mapStore, "db search")

	// validation should be done- email validate then save the user
	userObj.IsEmailValidated = true
	userObj.UpdatedAt = time.Now()

	uaString := mapStore["userAgent"].(string)
	// Parse() returns all attributes, including returning the full UA string last
	ua := uasurfer.Parse(uaString)
	log.Println(ua)

	_, err = userObj.Save(mapStore)
	if err != nil {
		log.Println("ERROR", err)

		log.Println("likely database down")

		acaGoMongoDBOpenWhiskUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusInternalServerError,
				"body":       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_INTERNAL", nil, language),
			},
			"text/html",
		)
		return
	}

	acaGoUtilities.PrintMilestone(mapStore, "db saving / update user")

	log.Println("function exec", (time.Now().UnixNano()-t1.UnixNano())/int64(time.Millisecond))

	body := acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "EMAIL_VALIDATION_SUCCESS_HTML", map[string]string{"BRANDING": acaGoConfiguration.BRANDING}, language)
	if ua.OS.Name == uasurfer.OSiOS {
		body = acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "EMAIL_VALIDATION_SUCCESS_HTML_IOS", map[string]string{"BRANDING": acaGoConfiguration.BRANDING, "BRANDING2": acaGoConfiguration.BRANDING2}, language)
	}

	// return success
	acaGoMongoDBOpenWhiskUtilities.DoResponse(
		map[string]interface{}{
			"statusCode": http.StatusOK,
			"body":       body,
		},
		"text/html",
	)
	return

	// main
}
