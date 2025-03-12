package repositories

import (
	"blog/internal/domain/users"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoUserRepository struct {
	db *mongo.Database
}

type UserDocument struct {
	ID        string `bson:"_id"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Username  string `bson:"username"`
	Email     string `bson:"email"`
	Password  string `bson:"password"`
	Bio       string `bson:"bio"`
	Avatar    string `bson:"avatar"`
}

func NewMongoUserRepository(db *mongo.Database) users.Repository {
	return &MongoUserRepository{db: db}
}

func (m *MongoUserRepository) Save(ctx context.Context, u users.User) error {
	timeoutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
	defer cancel()
	opt := &options.ReplaceOptions{}
	opt.SetUpsert(true)
	filters := bson.M{"_id": u.ID}
	_, err := m.db.Collection("users").ReplaceOne(timeoutCtx, filters, m.serialize(u), opt)
	return err
}

func (m *MongoUserRepository) FindByEmailOrUsername(ctx context.Context, identifier string) (*users.User, error) {
	timeoutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
	defer cancel()
	filter := bson.M{
		"$or": []bson.M{
			{"email": identifier},
			{"username": identifier},
		},
	}
	result := m.db.Collection("users").FindOne(timeoutCtx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var doc UserDocument
	err := result.Decode(&doc)
	if err != nil {
		return nil, err
	}
	return m.deserialize(doc), nil
}

func (m *MongoUserRepository) serialize(u users.User) *UserDocument {
	return &UserDocument{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Bio:       u.Bio,
	}
}

func (m *MongoUserRepository) deserialize(u UserDocument) *users.User {
	return &users.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Bio:       u.Bio,
	}
}
