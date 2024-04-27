package provider

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"stoo-kv/config"
	"strings"
)

type mongoKv struct {
	Key   string
	Value string
}
type MongoClient struct {
	client     *mongo.Client
	cfg        *config.Config
	ctx        context.Context
	collection *mongo.Collection
}

func NewMongoClient(ctx context.Context, config *config.Config) (*MongoClient, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Providers.Mongo.MongoUri))
	return &MongoClient{client: client,
		cfg:        config,
		collection: client.Database(config.Providers.Mongo.DatabaseName).Collection(config.Providers.Mongo.CollectionName)}, err
}

func (m *MongoClient) Set(key string, value any) error {
	filter := bson.D{{key, value}}
	document := bson.D{{"$set", bson.D{{key, value}}}}
	opts := options.Update().SetUpsert(true)
	_, err := m.collection.UpdateOne(m.ctx, filter, document, opts)
	return err
}

func (m *MongoClient) Get(key string) (string, error) {
	kv := &mongoKv{}
	err := m.collection.FindOne(m.ctx, bson.D{{"key", key}}).Decode(kv)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", nil
	}
	return kv.Value, err
}

func (m *MongoClient) Delete(key string) error {
	_, err := m.collection.DeleteMany(m.ctx, bson.D{{"key", key}})
	return err
}

//func (m *MongoClient) GetAll() (map[string]string, error) {
//	keyValues := make(map[string]string)
//	results, err := m.findAll()
//	if err != nil {
//		return validateError(err)
//	}
//
//	for k, v := range results {
//		keyValues[strings.Split(k, "::")[2]] = v
//	}
//	return keyValues, nil
//}

func (m *MongoClient) GetByNameSpaceAndProfile(namespace, profile string) (map[string]string, error) {
	keyValues := make(map[string]string)
	results, err := m.findAll()
	if err != nil {
		return validateError(err)
	}

	for k, v := range results {
		if strings.HasPrefix(k, fmt.Sprintf("%s::%s", namespace, profile)) {
			keyValues[strings.Split(k, "::")[2]] = v
		}
	}
	return keyValues, nil
}

func (m *MongoClient) findAll() (map[string]string, error) {
	keyValues := make(map[string]string)
	cursor, err := m.collection.Find(m.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err := cursor.All(m.ctx, &results); err != nil {
		return nil, err
	}
	for _, v := range results {
		keyValues[v["key"].(string)] = v["value"].(string)
	}
	return keyValues, nil
}

func validateError(err error) (map[string]string, error) {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return map[string]string{}, nil
	}
	return nil, err
}
