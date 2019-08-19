package storage

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Quard/poindexter/internal/reading_list"
	"github.com/getsentry/sentry-go"
)

const listCollectionName = "readinglist"

type MongoStorage struct {
	uri  string
	conn *mongo.Database
}

func NewMongoStorage(uri string) MongoStorage {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	return MongoStorage{
		uri:  uri,
		conn: client.Database("poindexter"),
	}
}

func (s MongoStorage) Add(record *reading_list.Record) error {
	listCollection := s.conn.Collection(listCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	doc, errConv := asNewDocument(record)
	if errConv != nil {
		sentry.CaptureException(errConv)
		return errConv
	}

	newDocument, err := listCollection.InsertOne(ctx, doc)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	record.ID = newDocument.InsertedID.(primitive.ObjectID)

	return nil
}

func (s MongoStorage) List(userID string, processor func(record *reading_list.Record)) error {
	listCollection := s.conn.Collection(listCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	userObjectID, errConv := primitive.ObjectIDFromHex(userID)
	if errConv != nil {
		sentry.CaptureException(errConv)
		return errors.New("authentication problem")
	}
	cursor, err := listCollection.Find(ctx, bson.M{"user_id": userObjectID})
	if err != nil {
		sentry.CaptureException(err)
		return errors.New("unable to get reading list")
	}
	defer cursor.Close(context.Background())

	for cursor.Next(ctx) {
		var record reading_list.Record
		if err := cursor.Decode(&record); err != nil {
			// sentry.CaptureException(err)
			// break
			log.Println(err)
			break
		} else {
			processor(&record)
		}
	}

	return cursor.Err()
}

func (s MongoStorage) MarkAsRead(userID string, ID string) (reading_list.Record, error) {
	var record reading_list.Record

	listCollection := s.conn.Collection(listCollectionName)
	userObjectID, errConvUserID := primitive.ObjectIDFromHex(userID)
	if errConvUserID != nil {
		sentry.CaptureException(errConvUserID)
		return record, errors.New("authentication problem")
	}
	objectID, errConvObjectID := primitive.ObjectIDFromHex(ID)
	if errConvObjectID != nil {
		sentry.CaptureException(errConvObjectID)
		return record, errors.New("getting record error")
	}

	filter := bson.M{"_id": objectID, "user_id": userObjectID}
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := listCollection.UpdateOne(
		ctx,
		filter,
		bson.M{"$set": bson.M{"is_read": true}},
	)
	if err != nil {
		sentry.CaptureException(err)
		return record, errors.New("unable to mark as read")
	}
	listItem := listCollection.FindOne(ctx, filter)
	if listItem.Err() != nil {
		sentry.CaptureException(listItem.Err())
		return record, errors.New("unable to retrieve record")
	}
	err = listItem.Decode(&record)
	if err != nil {
		sentry.CaptureException(err)
		return record, errors.New("unable to retrieve record")
	}

	return record, nil
}

func asNewDocument(val interface{}) (bson.M, error) {
	var document bson.M
	bytes, err := bson.Marshal(val)
	if err != nil {
		sentry.CaptureException(err)
		return document, err
	}
	if err := bson.Unmarshal(bytes, &document); err != nil {
		sentry.CaptureException(err)
		return document, err
	}

	delete(document, "_id")

	return document, nil
}
