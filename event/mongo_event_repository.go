package event

import (
	"context"
	"time"

	"github.com/somatom98/badges/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEventRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoEventRepository(db *mongo.Database) *MongoEventRepository {
	return &MongoEventRepository{
		db:   db,
		coll: db.Collection("events"),
	}
}

func (r *MongoEventRepository) GetEventsByUserID(ctx context.Context, uid string) ([]domain.Event, error) {
	events := make([]domain.Event, 0)

	objectID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return []domain.Event{}, err
	}

	opts := options.Find().SetSort(bson.M{"ts": -1})
	filter := bson.M{"uid": objectID}

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return events, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			return events, err
		}

		events = append(events, elem.toDomain())
	}

	if err := cur.Err(); err != nil {
		return events, err
	}

	return events, nil
}

func (r *MongoEventRepository) GetEventsByIDs(ctx context.Context, uids ...string) ([]domain.Event, error) {
	events := make([]domain.Event, 0)

	objectIDs := []primitive.ObjectID{}
	for _, uid := range uids {
		objectID, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			return events, err
		}

		objectIDs = append(objectIDs, objectID)
	}

	opts := options.Find().SetSort(bson.M{"ts": -1})
	filter := bson.M{
		"uid": bson.M{
			"$in": objectIDs,
		},
	}

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return events, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			return events, err
		}

		events = append(events, elem.toDomain())
	}

	if err := cur.Err(); err != nil {
		return events, err
	}

	return events, nil
}

func (r *MongoEventRepository) AddUserEvent(ctx context.Context, event domain.Event) error {
	mongoEvent, err := toMongo(event)
	if err != nil {
		return err
	}

	_, err = r.coll.InsertOne(ctx, mongoEvent)
	if err != nil {
		return err
	}

	return nil
}

type mongoEvent struct {
	ID   string             `bson:"id"`
	UID  primitive.ObjectID `bson:"uid"`
	Type domain.EventType   `bson:"type"`
	Date time.Time          `bson:"ts"`
}

func (e mongoEvent) toDomain() domain.Event {
	return domain.Event{
		ID:   e.ID,
		UID:  e.UID.Hex(),
		Type: e.Type,
		Date: e.Date,
	}
}

func toMongo(e domain.Event) (mongoEvent, error) {
	uid, err := primitive.ObjectIDFromHex(e.UID)
	if err != nil {
		return mongoEvent{}, err
	}

	return mongoEvent{
		ID:   e.ID,
		UID:  uid,
		Type: e.Type,
		Date: e.Date,
	}, nil
}
