// Created by Jonee Ryan Ty
// Copyright ACloudApp

// test.js
'use strict';

const acaNodejsConfiguration = require('../../../../../configuration/configuration');
const acaNodejsConstants = require('../../../../../configuration/constants');
const acaNodejsMongoDBConfiguration = require('../../../configuration/configuration');

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
  let stage = event.resource.substring(1); // /dev/test/test
  stage = stage.substring(0, stage.indexOf('/'));

  // let stage = event.requestContext.stage;
  console.log("stage:", stage);
	mapStore["stage"] = stage;

  // console.log("event.headers.Authorization:", event.headers.Authorization);
  let token = "";
  let payload;

  if (event.headers.Authorization && event.headers.Authorization.toUpperCase().startsWith("BEARER")) {
    token = event.headers.Authorization.substr(6).trim();
  }
  // console.log("token:", token);

  if (!token) {
    // require token
    let ret = {"success":0, "message":"jwt token required", "message_key":"ERROR_JWT_REQUIRED"};
    return {
      statusCode: 400,
      body: JSON.stringify(
        ret,
        null,
        2
      ),
    };
  }

  try {
    payload = jwt.verify(token, acaNodejsConfiguration.JWT_SECRET_STRING);
  } catch (e) {
    console.log(e);
  }
  // console.log("payload:", payload);

  await doDBConnect(mapStore);

  // TODO finish   


  return {
    statusCode: 200,
    body: JSON.stringify(
      {
        message: 'test',
        input: event,
      },
      null,
      2
    ),
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

