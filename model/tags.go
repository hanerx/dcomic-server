package model

type ComicTag struct {
	TagID string `json:"tag_id" bson:"tag_id"`
	Title string `json:"title" bson:"title"`
	Cover string `json:"cover" bson:"cover"`
}

