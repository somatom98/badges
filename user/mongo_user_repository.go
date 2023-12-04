package user

import (
	"context"

	"github.com/somatom98/badges/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		db:   db,
		coll: db.Collection("users"),
	}
}

func (r *MongoUserRepository) GetUserByID(ctx context.Context, uid string) (domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return domain.User{}, err
	}

	var elem mongoUser
	if err = r.coll.FindOne(ctx, bson.M{"_id": objectID}).Decode(&elem); err != nil {
		return domain.User{}, err
	}

	return elem.toDomain(), nil
}

func (r *MongoUserRepository) GetUsersByManagerID(ctx context.Context, managerID string) ([]domain.User, error) {
	users := []domain.User{}

	objectID, err := primitive.ObjectIDFromHex(managerID)
	if err != nil {
		return users, err
	}

	opts := options.Find().SetSort(bson.M{"mid": 1})
	filter := bson.M{"mid": objectID}

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return users, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem mongoUser
		err := cur.Decode(&elem)
		if err != nil {
			return users, err
		}

		users = append(users, elem.toDomain())
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	return users, nil
}

type mongoUser struct {
	ID        primitive.ObjectID  `bson:"_id"`
	ManagerID *primitive.ObjectID `bson:"mid"`
	Name      string              `bson:"name"`
}

func (u mongoUser) toDomain() domain.User {
	managerID := u.ManagerID.Hex()
	return domain.User{
		ID:        u.ID.Hex(),
		ManagerID: &managerID,
		Name:      u.Name,
	}
}
