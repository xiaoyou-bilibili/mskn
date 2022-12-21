package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoService(database string) *Service {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://192.168.1.10:8313"))
	if err != nil {
		panic(err)
	}
	return &Service{
		db: client.Database(database),
	}
}

type Service struct {
	db *mongo.Database
}

func (s *Service) Collection(collection string) *Collection {
	return &Collection{
		db:   s.db,
		col:  s.db.Collection(collection),
		find: options.Find(),
	}
}

type Collection struct {
	db   *mongo.Database
	col  *mongo.Collection
	find *options.FindOptions
}

func (c *Collection) getCtx() context.Context {
	return context.TODO()
}

func (c *Collection) InsertOne(data interface{}) error {
	_, err := c.col.InsertOne(c.getCtx(), data)
	return err
}

func (c *Collection) getFilter(filter *Filter) interface{} {
	if filter != nil {
		return filter.Res()
	}

	return bson.M{}
}

func (c *Collection) getUpdate(update *Update) interface{} {
	if update != nil {
		return update.Res()
	}

	return bson.M{}
}

func (c *Collection) FindOne(filter *Filter, data interface{}) error {
	res := c.col.FindOne(c.getCtx(), c.getFilter(filter))
	if res.Err() != nil {
		return res.Err()
	}
	return res.Decode(data)
}

func (c *Collection) Skip(num int64) *Collection {
	c.find.SetSkip(num)
	return c
}

func (c *Collection) Limit(limit int64) *Collection {
	c.find.SetLimit(limit)
	return c
}

func (c *Collection) FindMany(filter *Filter, data interface{}) error {
	cursor, err := c.col.Find(c.getCtx(), c.getFilter(filter), c.find)
	if err != nil {
		return err
	}
	return cursor.All(c.getCtx(), data)
}

func (c *Collection) FindByPage(filter *Filter, no, size int64, data interface{}) (int64, error) {
	if no <= 0 {
		no = 1
	}
	if size <= 0 {
		size = 20
	}
	option := options.Find()
	option.SetLimit(size).SetSkip((no - 1) * size).SetSort(bson.M{"_id": -1})
	// 获取数据总数
	total, err := c.col.CountDocuments(c.getCtx(), c.getFilter(filter))
	if err != nil {
		return 0, err
	}
	cursor, err := c.col.Find(c.getCtx(), c.getFilter(filter), option)
	if err != nil {
		return 0, err
	}
	return total, cursor.All(c.getCtx(), data)
}

func (c *Collection) Count(filter *Filter) (int64, error) {
	return c.col.CountDocuments(c.getCtx(), c.getFilter(filter))
}

func (c *Collection) UpdateOne(filter *Filter, update *Update) error {
	_, err := c.col.UpdateOne(c.getCtx(), c.getFilter(filter), c.getUpdate(update))
	return err
}

func (c *Collection) DeleteOne(filter *Filter) error {
	_, err := c.col.DeleteOne(c.getCtx(), c.getFilter(filter))
	return err
}
