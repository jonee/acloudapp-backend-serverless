// Created by Jonee Ryan Ty

package models

import ()

const GEO_JSON_TYPE_POINT = "Point"

// thanks http://icchan.github.io/2014/10/18/geospatial-querying-with-go-and-mongodb/
type GeoJson struct {
	Type        string    `json:"type,omitempty" bson:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
}
