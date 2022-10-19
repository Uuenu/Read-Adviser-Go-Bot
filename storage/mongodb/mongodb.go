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
	Ctx    context.Context
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
		Ctx:    context.TODO(),
	}

	return &s
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = lib.WrapIfErr("can't save page", err) }()

	userCollection := s.DB.Collection(page.UserName)
	fmt.Println(userCollection)
	if userCollection == nil {
		s.DB.CreateCollection(nil, page.UserName, nil)
	}
	doc := bson.D{{"page_url", page.URL}}
	result, err := userCollection.InsertOne(context.TODO(), doc)
	fmt.Println(result)
	return err

}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	userUrl, err := s.DB.Collection(p.UserName).Find(s.Ctx, bson.M{"page_url": p.URL})
	fmt.Println(userUrl, " ", err)
	// if err != nil {
	// 	return false, lib.Wrap("can't check data in db", err)
	// }
	// if userUrl.Current.String() != "" {
	// 	return false, nil
	// }
	// return true, nil
	return false, nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {

	// urlDocuments, err := s.DB.Collection(userName).Find(context.TODO(), bson.D{})
	// if err != nil {
	// 	return nil, lib.WrapIfErr("can't connect to db", err)
	// }

	// links := make([]string, 0)

	// for urlDocuments.Next(context.TODO()) {
	// 	var result bson.M
	// 	err := urlDocuments.Decode(&result)
	// 	// If there is a cursor.Decode error
	// 	if err != nil {
	// 		return nil, lib.WrapIfErr("cursor.Next() error:", err)

	// 		// If there are no cursor.Decode errors
	// 	} else {
	// 		links = append(links, fmt.Sprintf("%v", result["page_url"]))
	// 	}
	// }

	// rand.Seed(time.Now().UnixNano())
	// n := rand.Intn(len(links))

	// link := links[n]

	// var p storage.Page
	// p.URL = link

	//return &p, nil

	return nil, nil
}

func (s Storage) Remove(p *storage.Page) error {
	// _, err := s.DB.Collection(p.UserName).DeleteOne(s.Ctx, bson.M{"page": p.URL})
	// if err != nil {
	// 	fmt.Println("hello")
	// 	return lib.WrapIfErr("can'n remove page from collection", err)
	// }
	return nil
}
