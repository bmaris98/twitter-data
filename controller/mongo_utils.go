package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var database = "twitter_data"

type MongoContext struct {
	Ctx         context.Context
	Client      *mongo.Client
	Prompts     *mongo.Collection
	UnsafeStats *mongo.Collection
	Reports     *mongo.Collection
}

type Prompt struct {
	Query      string  `bson:"_id" json:"query"`
	IsActive   bool    `bson:"is_active" json:"isActive"`
	LastIdRead float64 `bson:"last_id_read" json:"lastIdRead"`
}

type Stat struct {
	Query     string `bson:"query" json:"query"`
	Timestamp uint64 `bson:"timestamp" json:"timestamp"`
	Value     uint64 `bson:"value" json:"value"`
}

type Report struct {
	Query     string `bson:"query" json:"query"`
	Data      string `bson:"data" json:"data"`
	Id        string `bson:"_id" json:"id"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
}

func (mongoCtx *MongoContext) InitMongo() {
	auth := "mongodb://mongoadmin:admin@twitter-mongo-host:27017/?serverSelectionTimeoutMS=5000&connectTimeoutMS=10000&authSource=admin&authMechanism=SCRAM-SHA-256"
	// auth := "mongodb://mongoadmin:admin@localhost:27017/?serverSelectionTimeoutMS=5000&connectTimeoutMS=10000&authSource=admin&authMechanism=SCRAM-SHA-256"
	clientOptions := options.Client().ApplyURI(auth)
	mongoCtx.Ctx = context.TODO()
	client, err := mongo.Connect(mongoCtx.Ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	mongoCtx.Client = client

	mongoCtx.Prompts = mongoCtx.Client.Database(database).Collection("prompts")
	mongoCtx.UnsafeStats = mongoCtx.Client.Database(database).Collection("unsafe_stats")
	mongoCtx.Reports = mongoCtx.Client.Database(database).Collection("reports")
}

func (mongoCtx MongoContext) AddPrompt(p Prompt) {
	_, err := mongoCtx.Client.Database(database).Collection("prompts").InsertOne(mongoCtx.Ctx, p)
	if err != nil {
		log.Println(err)
	}
}

func (mongoCtx MongoContext) TogglePromptStatus(prompt string) {
	originalPrompt := mongoCtx.FindOne(prompt)
	originalPrompt.IsActive = !originalPrompt.IsActive
	filter := bson.M{"_id": prompt}
	update := bson.M{"$set": bson.M{"is_active": originalPrompt.IsActive}}
	_, err := mongoCtx.Prompts.UpdateOne(mongoCtx.Ctx, filter, update)
	if err != nil {
		panic(err)
	}
}

func (mongoCtx MongoContext) UpdatePromptLastRead(prompt string, lastRead float64) {
	filter := bson.M{"_id": prompt}
	update := bson.M{"$set": bson.M{"last_id_read": lastRead}}
	_, err := mongoCtx.Prompts.UpdateOne(mongoCtx.Ctx, filter, update)
	if err != nil {
		panic(err)
	}
}

func (mongoCtx MongoContext) FindOne(id string) Prompt {
	filter := bson.M{"_id": id}
	var result Prompt
	response := mongoCtx.Prompts.FindOne(context.TODO(), filter)
	err := response.Err()
	if err != nil {
		log.Println(err)
	}
	err = response.Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Id did not match any entries from DB.")
		}
		log.Println(err)
		panic(err)
	}
	return result
}

func (mongoCtx MongoContext) GetAllPrompts() []Prompt {
	cursor, err := mongoCtx.Prompts.Find(mongoCtx.Ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var prompts []Prompt

	if err = cursor.All(context.TODO(), &prompts); err != nil {
		log.Fatal(err)
	}
	return prompts
}

func (mongoCtx MongoContext) InsertPrompt(p Prompt) {
	doc := bson.M{"_id": p.Query, "is_active": false, "last_id_read": 0}
	_, err := mongoCtx.Prompts.InsertOne(mongoCtx.Ctx, doc)

	if err != nil {
		panic(err)
	}
}

func (mongoCtx MongoContext) ReadAllUnsafeStats(query string) []Stat {
	filter := bson.M{"query": query}
	cursor, err := mongoCtx.UnsafeStats.Find(mongoCtx.Ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	var stats []Stat

	if err = cursor.All(context.TODO(), &stats); err != nil {
		log.Fatal(err)
	}
	return stats
}

func (mongoCtx MongoContext) InsertUnsafeStat(s Stat) {
	doc := bson.M{"query": s.Query, "timestamp": s.Timestamp, "value": s.Value}
	_, err := mongoCtx.UnsafeStats.InsertOne(mongoCtx.Ctx, doc)

	if err != nil {
		panic(err)
	}
}

func (mongoCtx MongoContext) ReadAllReports(query string) []Report {
	filter := bson.M{"query": query}
	cursor, err := mongoCtx.Reports.Find(mongoCtx.Ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	var reports []Report

	if err = cursor.All(context.TODO(), &reports); err != nil {
		log.Fatal(err)
	}
	return reports
}

func (mongoCtx MongoContext) InsertReport(r Report) {
	doc := bson.M{"_id": r.Id, "query": r.Query, "data": r.Data, "timestamp": r.Timestamp}
	_, err := mongoCtx.Reports.InsertOne(mongoCtx.Ctx, doc)

	if err != nil {
		panic(err)
	}
}
