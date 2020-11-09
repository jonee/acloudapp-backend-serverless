// Created by Jonee Ryan Ty
// Copyright ACloudApp

// test.js
'use strict';

const acaNodejsConfiguration = require('../../../../../../node.js/configuration/configuration');
const acaNodejsConstants = require('../../../../../../node.js/configuration/constants');
const acaNodejsMongoDBConfiguration = require('../../../../../../node.js/databases/mongodb/configuration/configuration');

const jwt = require("jsonwebtoken");

const Promise = require('bluebird');
const validator = require('validator');
const mongoose = require('mongoose');
// const UserModel = require('./model/User.js');

mongoose.Promise = Promise;

var mapStore = {};


async function test(event) {
  console.log(event);

  // extract the stage - dev or prod in the url path
  // eg https://service.us.apiconnect.ibmcloud.com/gws/apigateway/api/5a3bb96a02a41a8cc7aebf6ce0a231a89bb0864b056d29cda4c566ab3a38d839/dev/account/is_logged_in
  let tmpStage = event["__ow_headers"]["x-forwarded-url"];
  tmpStage = tmpStage.substring(tmpStage.indexOf(event["__ow_headers"]["x-forwarded-prefix"]) + event["__ow_headers"]["x-forwarded-prefix"].length + 1); // eg api/5a3bb96a02a41a8cc7aebf6ce0a231a89bb0864b056d29cda4c566ab3a38d839/dev/account/is_logged_in
  let explode = tmpStage.split("/");
  let stage = explode[2];
  console.log("stage:", stage);
	mapStore["stage"] = stage;

  // console.log('event["__ow_headers"]["authorization"]:', event["__ow_headers"]["authorization"]);
  let token = "";
  let payload;

  if (event["__ow_headers"]["authorization"] && event["__ow_headers"]["authorization"].toUpperCase().startsWith("BEARER")) {
    token = event["__ow_headers"]["authorization"].substr(6).trim();
  }
  // console.log("token:", token);

  if (!token) { // require token
    return {
      body: {"success":0, "message":"jwt token required", "message_key":"ERROR_JWT_REQUIRED"},
      statusCode: 400,
      headers:{ 'Content-Type': 'application/json'}
    };
  }

  try {
    payload = jwt.verify(token, acaNodejsConfiguration.JWT_SECRET_STRING);
  } catch (e) {
    console.log(e);
  }
  console.log("payload:", payload);

  await doDBConnect(mapStore);

  // TODO finish


  return {
    body: {payload: `Hello world`},
    statusCode: 200,
    headers:{ 'Content-Type': 'application/json'}
  };

}


async function doDBConnect(mapStore) {
  let stage = mapStore["stage"];
  console.log("doDBConnect stage:", stage);

  // db config
	let mongoDatabase = acaNodejsMongoDBConfiguration.mongoConfig[stage]["MONGODB_DATABASE"];
	let mongoAuthDatabase = acaNodejsMongoDBConfiguration.mongoConfig[stage]["MONGODB_AUTH_DATABASE"];
	let mongoUser = acaNodejsMongoDBConfiguration.mongoConfig[stage]["MONGODB_REGULAR_USER"];
	let mongoUserPassword = acaNodejsMongoDBConfiguration.mongoConfig[stage]["MONGODB_REGULAR_USER_PASSWORD"];
	let mongoClusterSrvUri = acaNodejsMongoDBConfiguration.mongoConfig[stage]["MONGODB_CLUSTER_SRV_URI"];
	let mongoClient = acaNodejsMongoDBConfiguration.mongoConfig[stage]["mongoClient"];

	mapStore["mongoDatabase"] = mongoDatabase;

  // const mongoConnStr = "mongodb+srv://" + mongoUser + ":" + mongoUserPassword + "@" + mongoClusterSrvUri + "/" + mongoAuthDatabase + "?retryWrites=true&w=majority";
  const mongoConnStr = "mongodb+srv://" + mongoUser + ":" + mongoUserPassword + "@" + mongoClusterSrvUri + "/" + mongoDatabase + "?retryWrites=true&w=majority";

  if (!mongoClient || (mongoClient.readyState != 1 && mongoClient.readyState != 2)) { // 0: disconnected, 1: connected, 2: connecting, 3: disconnecting
    mongoClient = mongoose.createConnection(mongoConnStr, {useNewUrlParser: true, useUnifiedTopology: true});
    console.log("mongoClient connecting");

    // TODO schema stuffs here


    acaNodejsMongoDBConfiguration.mongoConfig[stage]["mongoClient"] = mongoClient; // save back
	  mapStore["mongoClient"] = mongoClient;

  } else {
    console.log("reusing mongoClient connection");
  }

  // TODO finish


} // doDBConnect


module.exports = {
    test,
};
