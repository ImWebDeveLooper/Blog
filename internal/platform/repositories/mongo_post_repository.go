package repositories

import (
	"blog/internal/domain/posts"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoPostRepository struct {
	db *mongo.Database
}

type PostDocument struct {
	ID      string `bson:"_id"`
	Title   string `bson:"title"`
	Content string `bson:"content"`
	Author  string `bson:"author"`
}

func NewMongoPostRepository(db *mongo.Database) posts.Repository {
	return &MongoPostRepository{db: db}
}

func (m *MongoPostRepository) Save(ctx context.Context, p posts.Post) error {
	timeoutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
	defer cancel()
	opt := &options.ReplaceOptions{}
	opt.SetUpsert(true)
	filters := bson.M{"_id": p.ID}
	_, err := m.db.Collection("posts").ReplaceOne(timeoutCtx, filters, m.serialize(p), opt)
	return err
}

//func (m *MongoPostRepository) Update(ctx context.Context, p posts.Post) error {
//	filter := bson.M{"author": p.Author}
//	updateResult, err := m.db.Collection("posts").ReplaceOne(ctx, filter, m.serialize(p))
//	if err != nil {
//		return err
//	}
//	if updateResult.MatchedCount == 0 {
//		return errors.New("no document found with the specified title")
//	}
//	return nil
//}

func (m *MongoPostRepository) serialize(p posts.Post) *PostDocument {
	return &PostDocument{
		ID:      p.ID,
		Title:   p.Title,
		Content: p.Content,
		Author:  p.Author,
	}
}

func (m *MongoPostRepository) deserialize(p PostDocument) *posts.Post {
	return &posts.Post{
		ID:      p.ID,
		Title:   p.Title,
		Content: p.Content,
		Author:  p.Author,
	}
}
