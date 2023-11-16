package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TodoBody struct {
	Title     string `json:"title,omitempty" bson:"title"`
	Completed bool   `json:"completed" bson:"completed"`
}

type TodoId struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"`
	User string             `json:"user,omitempty" bson:"user"`
}

type Todo struct {
	Title        string    `json:"title,omitempty" bson:"title,omitempty"`
	Completed    bool      `json:"completed" bson:"completed"`
	User         string    `json:"user,omitempty" bson:"user,omitempty"`
	Created_time time.Time `json:"created_time" bson:"created_time,omitempty"`
	Updated_time time.Time `json:"updated_time" bson:"updated_time,omitempty"`
	Id           string    `json:"id,omitempty" bson:"_id,omitempty"`
}

func (t *Todo) FetchFromBody(todoBody TodoBody) {
	t.Title = todoBody.Title
	t.Completed = todoBody.Completed
	t.Updated_time = time.Now()
}
