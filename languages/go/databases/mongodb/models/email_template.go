// Created by Jonee Ryan Ty
// Copyright ACloudApp

/**
 * EmailTemplate model class
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

// EmailTemplate class or struct definition
type EmailTemplate struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Language string `json:"language,omitempty" bson:"language,omitempty"`

	Subject string `json:"subject,omitempty" bson:"subject,omitempty"`

	TemplateHTML string `json:"template_html,omitempty" bson:"template_html,omitempty"`
	TemplateText string `json:"template_text,omitempty" bson:"template_text,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// EmailTemplate class Save function
func (r *EmailTemplate) Save(mapStore map[string]interface{}) (*EmailTemplate, error) {
	emailTemplateCol := mapStore["emailTemplateCol"].(*mongo.Collection)
	myCol := emailTemplateCol

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
func (r *EmailTemplate) ExportArrayPrivate() map[string]interface{} {
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
	if r.Subject != "" {
		val["subject"] = r.Subject
	}
	if r.TemplateHTML != "" {
		val["template_html"] = r.TemplateHTML
	}
	if r.TemplateText != "" {
		val["template_text"] = r.TemplateText
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
func (r *EmailTemplate) ExportArrayPublic() map[string]interface{} {
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
