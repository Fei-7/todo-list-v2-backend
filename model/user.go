package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Name     string             `bson:"name" json:"name"`
	Detail   string             `bson:"detail" json:"detail"`
	Time     int64              `bson:"time" json:"time"`
	Priority int                `bson:"priority" json:"priority"`
	Tags     []string           `bson:"tags" json:"tags"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password []byte             `bson:"password" json:"-"`
	Tasks    []Task             `bson:"tasks" json:"tasks"`
}
