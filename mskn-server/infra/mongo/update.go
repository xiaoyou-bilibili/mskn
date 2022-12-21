package mongo

import "go.mongodb.org/mongo-driver/bson"

func NewUpdate() *Update {
	return &Update{
		data: bson.M{},
	}
}

type Update struct {
	data bson.M
}

func (u *Update) SetField(key string, value interface{}) *Update {
	if _, ok := u.data["$set"]; !ok {
		u.data["$set"] = bson.M{}
	}
	u.data["$set"].(bson.M)[key] = value
	return u
}

func (u *Update) Set(data interface{}) *Update {
	u.data["$set"] = data
	return u
}

func (u *Update) Res() interface{} {
	return u.data
}
