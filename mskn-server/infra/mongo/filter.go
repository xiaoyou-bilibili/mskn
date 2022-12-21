package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewFilter() *Filter {
	return &Filter{
		condition: bson.M{},
	}
}

type Filter struct {
	condition bson.M
}

func (f *Filter) fieldCheck(field string) {
	if _, ok := f.condition[field]; !ok {
		f.condition[field] = bson.M{}
	}
}

func (f *Filter) Id(id string) *Filter {
	oid, _ := primitive.ObjectIDFromHex(id)
	f.condition["_id"] = oid
	return f
}

func (f *Filter) FieldEq(key string, value string) *Filter {
	f.condition[key] = value
	return f
}

func (f *Filter) FieldIn(key string, value interface{}) *Filter {
	f.fieldCheck(key)
	f.condition[key].(bson.M)["$in"] = value
	return f
}

func (f *Filter) Res() interface{} {
	return f.condition
}
