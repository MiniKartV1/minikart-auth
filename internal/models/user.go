package models

import (
	"time"

	_ "go.mongodb.org/mongo-driver/bson"
)

type User struct {
	FirstName    string     `bson:"firstname,omitempty"`
	LastName     string     `bson:"lastname,omitempty"`
	Email        string     `bson:"email,omitempty"`
	Password     string     `json:"-" bson:"password,omitempty"`
	LastSignedIn *time.Time `bson:"lastSignedIn,omitempty"` // Pointer to time.Time allows for a nil value
	Verified     bool       `bson:"verified,omitempty"`     // Defaults to false
	IsActive     bool       `bson:"isActive,omitempty"`
}
