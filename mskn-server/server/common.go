package server

import "mskn-server/infra/mongo"

var mongoServer = mongo.NewMongoService("mskn")
var CodeCollection = mongoServer.Collection("code")
var TaskCollection = mongoServer.Collection("task")
var NodeCollection = mongoServer.Collection("node")
var DataCollection = mongoServer.Collection("data")
var RecordCollection = mongoServer.Collection("record")
