// Created by Jonee Ryan Ty
// Copyright ACloudApp

/**
 * LoginLog model class
 */

package models

import (
	"context"
	// "errors"
	// "log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

// LoginLog class or struct definition
type LoginLog struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	UserId primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`

	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`

	Version         string    `json:"version,omitempty" bson:"version,omitempty"`
	ApplicationType string    `json:"application_type,omitempty" bson:"application_type,omitempty"`
	Expiry          time.Time `json:"expiry,omitempty" bson:"expiry,omitempty"`
	IsValid         bool      `json:"is_valid" bson:"is_valid"`
	InvalidReason   string    `json:"invalid_reason,omitempty" bson:"invalid_reason,omitempty"`
	Secret          string    `json:"secret,omitempty" bson:"secret,omitempty"`

	LoggedOutAt time.Time `json:"logged_out_at,omitempty" bson:"logged_out_at,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// LoginLog class Save function
func (r *LoginLog) Save(mapStore map[string]interface{}) (*LoginLog, error) {
	loginLogCol := mapStore["loginLogCol"].(*mongo.Collection)
	myCol := loginLogCol

	var err error
	ctx := context.Background()

	if r.Id.IsZero() {
		r.Id = primitive.NewObjectID()
		_, err = myCol.InsertOne(ctx, &r)

	} else {
		_, err = myCol.UpdateOne(ctx, bson.M{"_id": r.Id}, bson.M{"$set": r.ExportArrayPrivate()})
	}

	return r, err

	// Save
}

// exportarrayprivate is saveable
func (r *LoginLog) ExportArrayPrivate() map[string]interface{} {
	val := make(map[string]interface{})

	if !r.Id.IsZero() {
		val["_id"] = r.Id
	}

	if !r.UserId.IsZero() {
		val["user_id"] = r.UserId
	}

	if r.Username != "" {
		val["username"] = r.Username
	}
	if r.Email != "" {
		val["email"] = r.Email
	}

	if r.Version != "" {
		val["version"] = r.Version
	}
	if r.ApplicationType != "" {
		val["application_type"] = r.ApplicationType
	}

	if !r.Expiry.IsZero() {
		val["expiry"] = r.Expiry
	}

	val["is_valid"] = r.IsValid

	if r.InvalidReason != "" {
		val["invalid_reason"] = r.InvalidReason
	}
	if r.Secret != "" {
		val["secret"] = r.Secret
	}

	if !r.LoggedOutAt.IsZero() {
		val["logged_out_at"] = r.LoggedOutAt
	}

	if !r.CreatedAt.IsZero() {
		val["created_at"] = r.CreatedAt
	}

	if !r.UpdatedAt.IsZero() {
		val["updated_at"] = r.UpdatedAt
	}

	return val

	// ExportArrayPrivate
}

// public does not have the passwords and anything sensitive / secure / private, also dates are ints
func (r *LoginLog) ExportArrayPublic() map[string]interface{} {
	val := r.ExportArrayPrivate()

	if !r.Id.IsZero() {
		val["_id"] = r.Id.Hex()
	}

	if !r.UserId.IsZero() {
		val["user_id"] = r.UserId.Hex()
	}

	// dates to ints
	if !r.Expiry.IsZero() {
		val["expiry"] = r.Expiry.Unix()
	}

	if !r.LoggedOutAt.IsZero() {
		val["logged_out_at"] = r.LoggedOutAt.Unix()
	}

	if !r.CreatedAt.IsZero() {
		val["created_at"] = r.CreatedAt.Unix()
	}

	if !r.UpdatedAt.IsZero() {
		val["updated_at"] = r.UpdatedAt.Unix()
	}

	return val

	// ExportArrayPublic
}
