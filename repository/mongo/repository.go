package mongo

import (
	"context"
	"time"

	"github.com/MrWormHole/url-shortener/shortener"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout))
	defer cancel()

	client, err := mongo.Connect(context, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, err
}

func NewMongoRepository(mongoURL string, mongoDB string, mongoTimeout int) (shortener.RedirectRepository, error) {
	repository := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepository")
	}
	repository.client = client

	return repository, nil
}

func (r *mongoRepository) Find(hash string) (*shortener.Redirect, error) {
	context, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	redirect := &shortener.Redirect{}
	collection := r.client.Database(r.database).Collection("redirects")
	filter := bson.M{"hash": hash}

	err := collection.FindOne(context, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	return redirect, nil
}

func (r *mongoRepository) Store(redirect *shortener.Redirect) error {
	context, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("redirects")

	_, err := collection.InsertOne(
		context,
		bson.M{
			"hash":       redirect.Hash,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
