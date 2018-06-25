package model

import "github.com/globalsign/mgo/bson"

type (
	Story struct {
		ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Author  string        `json:"author" bson:"author"`
		Title   string        `json:"title" bson:"title"`
		Content string        `json:"content" bson:"content"`
	}
)
