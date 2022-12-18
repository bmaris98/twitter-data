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
	Ctx     context.Context
	Client  *mongo.Client
	Prompts *mongo.Collection
}

type Prompt struct {
	Query      string  `bson:"_id" json:"query"`
	IsActive   bool    `bson:"is_active" json:"isActive"`
	LastIdRead float64 `bson:"last_id_read" json:"lastIdRead"`
}

func (mongoCtx *MongoContext) InitMongo() {
	auth := "mongodb://mongoadmin:admin@twitter-mongo-host:27017/?serverSelectionTimeoutMS=5000&connectTimeoutMS=10000&authSource=admin&authMechanism=SCRAM-SHA-256"
	clientOptions := options.Client().ApplyURI(auth)
	mongoCtx.Ctx = context.TODO()
	client, err := mongo.Connect(mongoCtx.Ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	mongoCtx.Client = client

	mongoCtx.Prompts = mongoCtx.Client.Database(database).Collection("prompts")
	log.Println(mongoCtx.Prompts)
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
	log.Println(id)
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
