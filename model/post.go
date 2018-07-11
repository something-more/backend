package model

import "github.com/globalsign/mgo/bson"

type (
	Post struct {
		ID             bson.ObjectId `json:"id" bson:"_id,omitempty"`
		AuthorID       bson.ObjectId `json:"author_id" bson:"author_id"`
		AuthorNickname string        `json:"author_nickname" bson:"author_nickname"`
		Thumbnail      string        `json:"thumbnail" bson:"thumbnail"`
		DateCreated    string        `json:"date_created" bson:"date_created"`
		DateModified   string        `json:"date_modified" bson:"date_modified"`
		Title          string        `json:"title" bson:"title"`
		Content        string        `json:"content" bson:"content"`
		IsPublished    bool          `json:"is_published" bson:"is_published"`
		Category       string           `json:"category" bson:"category"`
	}
)
