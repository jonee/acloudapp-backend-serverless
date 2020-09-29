// Created by Jonee Ryan Ty
// Copyright ACloudApp

/**
 * Translation model class
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

// Translation class or struct definition
type Translation struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Language string `json:"language,omitempty" bson:"language,omitempty"`

	Translation string `json:"translation,omitempty" bson:"translation,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// Translation class Save function
func (r *Translation) Save(mapStore map[string]interface{}) (*Translation, error) {
	translationCol := mapStore["translationCol"].(*mongo.Collection)
	myCol := translationCol

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
func (r *Translation) ExportArrayPrivate() map[string]interface{} {
	val := make(map[string]interface{})

	if !r.Id.IsZero() {
		val["_id"] = r.Id
	}

	if r.Name != "" {
		val["name"] = r.Name
	}
	if r.Language != "" {
		val["language"] = r.Language
	}
	if r.Translation != "" {
		val["translation"] = r.Translation
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
func (r *Translation) ExportArrayPublic() map[string]interface{} {
	val := r.ExportArrayPrivate()

	if !r.Id.IsZero() {
		val["_id"] = r.Id.Hex()
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
