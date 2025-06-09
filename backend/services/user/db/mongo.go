package db

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type User struct {
	Id       bson.ObjectID `bson:"_id,omitempty"`
	Name     string        `bson:"name,omitempty"`
	Age      int32         `bson:"age,omitempty"`
	Greeting string        `bson:"greeting,omitempty"`
	Salary   int32         `bson:"salary,omitempty"`
	Power    string        `bson:"power,omitempty"`
}

var _ = loadLocalEnv()
var (
	db   = GetEnv("MONGO_DATABASE")
	user = GetEnv("MONGO_USER")
	pwd  = GetEnv("MONGO_PWD")
	coll = GetEnv("MONGO_COLLECTION")
	addr = GetEnv("MONGO_ADDR")
)

var MongoClient *mongo.Client

func NewClient(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(addr).
		SetAuth(options.Credential{
			AuthSource: db,
			Username:   user,
			Password:   pwd,
		}))
	if err != nil {
		return nil, errors.New("invalid mongodb options")
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.New("cannot connect to mongodb instance")
	}
	return client, nil
}

func UpsertOne(ctx context.Context, user *User) error {
	collection := MongoClient.Database(db).Collection(coll)

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": bson.M{"age": user.Age, "name": user.Name,
		"salary": user.Salary, "greeting": user.Greeting, "power": user.Power}}

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func FindOne(ctx context.Context, id bson.ObjectID) (*User, error) {
	collection := MongoClient.Database(db).Collection(coll)

	var data User

	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func Find(ctx context.Context) (*[]User, error) {
	collection := MongoClient.Database(db).Collection(coll)

	var data []User
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		data = append(data, user)
	}

	return &data, nil
}

func loadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("MONGO_ADDR"); !runningInContainer {
		err := godotenv.Load("../.env.local")
		if err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}

func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln("Environment variable not found: ", key)
	}
	return value
}
