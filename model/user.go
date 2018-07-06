package model

import "github.com/globalsign/mgo/bson"

type (
	User struct {
		ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Email    string        `json:"email" bson:"email,omitempty"`
		Nickname string        `json:"nickname" bson:"nickname,omitempty"`
		Password string        `json:"password" bson:"password,omitempty"`
		Token    string        `json:"token,omitempty" bson:"-"`
		IsActive bool          `json:"is_active" bson:"is_active"`
		IsAdmin  bool          `json:"is_admin" bson:"is_admin"`
		IsStaff  bool          `json:"is_staff" bson:"is_staff"`
	}
)
