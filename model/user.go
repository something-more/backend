package model

import "github.com/globalsign/mgo/bson"

type (
	User struct {
		ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Email    string        `json:"email" bson:"email,omitempty"`
		Password string        `json:"password" bson:"password,omitempty"`
		Token    string        `json:"token,omitempty" bson:"-"`
		IsActive bool          `json:"is_active" bson:"is_active"`
	}
)
