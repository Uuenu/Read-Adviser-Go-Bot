package mongodb

import (
	"context"
	"fmt"
	lib "telegram-bot-go/lib/e"
	"telegram-bot-go/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbURI = "mongodb://localhost:27017"
)

type Storage struct {
	// mongo db client
	Client *mongo.Client
	DB     *mongo.Database
}

func New() *Storage {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil
	}
	database := client.Database("Read-Adviser-Bot")
	if err != nil {
		return nil
	}
	s := Storage{
		Client: client,
		DB:     database,
	}

	return &s
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = lib.WrapIfErr("can't save page", err) }()

	userCollection := s.DB.Collection(page.UserName)
	fmt.Println(userCollection)
	if userCollection == nil {
		fmt.Println(111)
		s.DB.CreateCollection(nil, page.UserName, nil)
	}
	doc := bson.D{{"page_url", page.URL}}
	result, err := userCollection.InsertOne(context.TODO(), doc)
	fmt.Println(result)
	return err

}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	return false, nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	return nil, nil
}

func (s Storage) Remove(p *storage.Page) error {
	return nil
}
