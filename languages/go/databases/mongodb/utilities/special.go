// Created by Jonee Ryan Ty
// Copyright ACloudApp

/**
 * Special Utility functions
 */

package utilities

import (
	"context"
	"log"
	"strings"
	"time"

	// acaGoConfiguration "acloudapp.org/configuration"
	acaGoMongoDBConfiguration "acloudapp.org/databases/mongodb/configuration"
	acaGoMongoDBModels "acloudapp.org/databases/mongodb/models"
	// acaGoUtilities "acloudapp.org/utilities"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

// connect to db
func DoDBConnect(mapStore map[string]interface{}) {
	var err error
	var tm time.Time

	stage := mapStore["stage"].(string)

	// db config
	mongoDatabase := acaGoMongoDBConfiguration.MongoConfig[stage]["MONGODB_DATABASE"].(string)
	mongoAuthDatabase := acaGoMongoDBConfiguration.MongoConfig[stage]["MONGODB_AUTH_DATABASE"].(string)
	mongoUser := acaGoMongoDBConfiguration.MongoConfig[stage]["MONGODB_REGULAR_USER"].(string)
	mongoUserPassword := acaGoMongoDBConfiguration.MongoConfig[stage]["MONGODB_REGULAR_USER_PASSWORD"].(string)
	mongoClusterSrvUri := acaGoMongoDBConfiguration.MongoConfig[stage]["MONGODB_CLUSTER_SRV_URI"].(string)
	mongoClient, _ := acaGoMongoDBConfiguration.MongoConfig[stage]["mongoClient"].(*mongo.Client)

	mapStore["mongoDatabase"] = mongoDatabase

	// connect mongodb if needed or reuse
	connectMongo := false
	hasPingSuccess := false
	if mongoClient == nil {
		connectMongo = true
		log.Println("mongoClient is nil")
	} else {
		// Check the connection
		tm = time.Now()
		err = mongoClient.Ping(context.TODO(), nil)
		log.Println("ping 1", (time.Now().UnixNano()-tm.UnixNano())/int64(time.Millisecond))

		if err != nil {
			connectMongo = true
			log.Println("mongoClient connected but not pinging")
		} else {
			hasPingSuccess = true
		}
	}

	if connectMongo {
		tm = time.Now()

		// Set client options
		clientOptions := options.Client().ApplyURI("mongodb+srv://" + mongoUser + ":" + mongoUserPassword + "@" + mongoClusterSrvUri + "/" + mongoAuthDatabase + "?retryWrites=true&w=majority")

		// Connect to MongoDB
		mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
		// defer mongoClient.Disconnect(context.TODO()) // don't close!

		log.Println("mongo connect", (time.Now().UnixNano()-tm.UnixNano())/int64(time.Millisecond))
		log.Println("new mongodb connection")

	} else {
		log.Println("re using mongodb")
	}

	if !hasPingSuccess {
		tm = time.Now()
		err = mongoClient.Ping(context.TODO(), nil)
		log.Println("ping 2", (time.Now().UnixNano()-tm.UnixNano())/int64(time.Millisecond))
	}

	acaGoMongoDBConfiguration.MongoConfig[stage]["mongoClient"] = mongoClient
	mapStore["mongoClient"] = mongoClient

	emailTemplateCol := mongoClient.Database(mongoDatabase).Collection("email_template")
	mapStore["emailTemplateCol"] = emailTemplateCol

	loginLogCol := mongoClient.Database(mongoDatabase).Collection("login_log")
	mapStore["loginLogCol"] = loginLogCol

	translationCol := mongoClient.Database(mongoDatabase).Collection("translation")
	mapStore["translationCol"] = translationCol

	userCol := mongoClient.Database(mongoDatabase).Collection("user")
	mapStore["userCol"] = userCol

	// DoDBConnect
}

/*
func AttachFollowStatus(jwtUserObj models.User, followObj models.Follow) (bool, bool, bool) {
	userId := jwtUserObj.Id
	userIdHex := userId.Hex()
	rec := followObj

	isFollowing := false
	isFollower := false
	isFriend := false

	// side := false // left false, right true where this jwtUser is
	// otherUserIdHex := ""
	// otherUsername := ""
	if userIdHex == rec.LeftUserId.Hex() {
		// side = false
		// otherUserIdHex = rec.RightUserId.Hex()
		// otherUsername = rec.RightUsername
		isFollowing = rec.LeftValue
		isFollower = rec.RightValue

	} else {
		// side = true
		// otherUserIdHex = rec.LeftUserId.Hex()
		// otherUsername = rec.LeftUsername
		isFollowing = rec.RightValue
		isFollower = rec.LeftValue
	}

	isFriend = isFollowing && isFollower

	return isFollowing, isFollower, isFriend

	// AttachFollowStatus
}
*/

func GetBackendTranslation(mapStore map[string]interface{}, messageKey string, messageParameters map[string]string, language string) string {
	ret := messageKey

	filter := bson.M{"language": language, "name": messageKey}

	ctx := context.Background()
	translationCol := mapStore["translationCol"].(*mongo.Collection)

	var translationObj acaGoMongoDBModels.Translation
	err := translationCol.FindOne(ctx, filter).Decode(&translationObj)
	if err == nil { // object found
		ret = translationObj.Translation
		if messageParameters != nil { // has parameters
			for k, v := range messageParameters {
				ret = strings.Replace(ret, "%%"+k+"%%", v, -1)
			}
		}

	} else {
		// log.Println(err)
		log.Println("translation not found", messageKey, language)
	}

	return ret

	// GetBackendTranslation
}
