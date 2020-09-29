// Created by Jonee Ryan Ty
// Copyright ACloudApp

/**
 * User model class
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

// User class or struct definition
type User struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`

	PasswordHash            string    `json:"password_hash,omitempty" bson:"password_hash,omitempty"`
	PasswordSalt            string    `json:"password_salt,omitempty" bson:"password_salt,omitempty"`
	PasswordTemporaryHash   string    `json:"password_temporary_hash,omitempty" bson:"password_temporary_hash,omitempty"`
	PasswordTemporaryExpiry time.Time `json:"password_temporary_expiry,omitempty" bson:"password_temporary_expiry,omitempty"`

	IsEmailValidated bool   `json:"is_email_validated" bson:"is_email_validated"`
	ValidationSecret string `json:"validation_secret,omitempty" bson:"validation_secret,omitempty"`

	LoginCount int `json:"login_count,omitempty" bson:"login_count,omitempty"`

	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	Bio     string `json:"bio,omitempty" bson:"bio,omitempty"`
	Town    string `json:"town,omitempty" bson:"town,omitempty"`
	Website string `json:"website,omitempty" bson:"website,omitempty"`
	Phone   string `json:"phone,omitempty" bson:"phone,omitempty"`
	Gender  string `json:"gender,omitempty" bson:"gender,omitempty"` // M or F

	Latitude  float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`   // -90 to +90
	Longitude float64 `json:"longitude,omitempty" bson:"longitude,omitempty"` // -180 to 180
	Location  GeoJson `json:"location,omitempty" bson:"location,omitempty"`

	IsBlocked bool `json:"is_blocked" bson:"is_blocked"`

	Access string `json:"access,omitempty" bson:"access,omitempty"` // C for customer

	Language string `json:"language,omitempty" bson:"language,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// User class Save function
func (r *User) Save(mapStore map[string]interface{}) (*User, error) {
	userCol := mapStore["userCol"].(*mongo.Collection)
	myCol := userCol

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
func (r *User) ExportArrayPrivate() map[string]interface{} {
	val := make(map[string]interface{})

	if !r.Id.IsZero() {
		val["_id"] = r.Id
	}

	if r.Username != "" {
		val["username"] = r.Username
	}
	if r.Email != "" {
		val["email"] = r.Email
	}

	if r.PasswordHash != "" {
		val["password_hash"] = r.PasswordHash
	}
	if r.PasswordSalt != "" {
		val["password_salt"] = r.PasswordSalt
	}
	if r.PasswordTemporaryHash != "" {
		val["password_temporary_hash"] = r.PasswordTemporaryHash
	}

	val["is_email_validated"] = r.IsEmailValidated

	if r.ValidationSecret != "" {
		val["validation_secret"] = r.ValidationSecret
	}

	if r.LoginCount != 0 {
		val["login_count"] = r.LoginCount
	}

	if r.Name != "" {
		val["name"] = r.Name
	}
	if r.Bio != "" {
		val["bio"] = r.Bio
	}
	if r.Town != "" {
		val["town"] = r.Town
	}
	if r.Website != "" {
		val["website"] = r.Website
	}
	if r.Phone != "" {
		val["phone"] = r.Phone
	}
	if r.Gender != "" {
		val["gender"] = r.Gender
	}

	val["is_blocked"] = r.IsBlocked

	if r.Access != "" {
		val["access"] = r.Access
	}

	if r.Language != "" {
		val["language"] = r.Language
	}

	if !r.PasswordTemporaryExpiry.IsZero() {
		val["password_temporary_expiry"] = r.PasswordTemporaryExpiry
	}

	if r.Latitude != 0 || r.Longitude != 0 { // let us not allow only (0, 0)
		val["latitude"] = r.Latitude
		val["longitude"] = r.Longitude

		var gj GeoJson
		gj.Type = GEO_JSON_TYPE_POINT
		gj.Coordinates = []float64{r.Longitude, r.Latitude}

		val["location"] = gj
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
func (r *User) ExportArrayPublic() map[string]interface{} {
	val := r.ExportArrayPrivate()

	if !r.Id.IsZero() {
		val["_id"] = r.Id.Hex()
	}

	// do not show senstive / secure / private stuffs
	delete(val, "password_hash")
	delete(val, "password_salt")
	delete(val, "password_temporary_hash")
	delete(val, "password_temporary_expiry")

	delete(val, "validation_secret")

	delete(val, "location")

	// dates to ints
	if !r.PasswordTemporaryExpiry.IsZero() {
		val["password_temporary_expiry"] = r.PasswordTemporaryExpiry.Unix()
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
