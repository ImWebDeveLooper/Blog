package repositories

import (
	"blog/internal/domain/posts"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// MongoPostRepository manages post data persistence using MongoDB.
type MongoPostRepository struct {
	db *mongo.Database
}

// PostDocument represents a post as stored in MongoDB.
type PostDocument struct {
	ID          string       `bson:"_id"`
	Title       string       `bson:"title"`
	Content     string       `bson:"content"`
	Author      string       `bson:"author"`
	Slug        string       `bson:"slug"`
	Status      posts.Status `bson:"status"`
	PublishedAt time.Time    `bson:"published_at"`
	CreatedAt   time.Time    `bson:"created_at"`
	UpdatedAt   time.Time    `bson:"updated_at"`
}

// NewMongoPostRepository creates a new MongoPostRepository with the given MongoDB database.
func NewMongoPostRepository(db *mongo.Database) posts.Repository {
	return &MongoPostRepository{db: db}
}

// Save creates a new post in MongoDB, ensuring no duplicate ID exists.
func (m *MongoPostRepository) Save(ctx context.Context, p posts.Post) error {
	timeoutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
	defer cancel()
	_, err := m.db.Collection("posts").InsertOne(timeoutCtx, m.serialize(p))
	return err
}

// FindByID retrieves a post by its ID from MongoDB.
func (m *MongoPostRepository) FindByID(ctx context.Context, id string) (*posts.Post, error) {
	timeoutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
	defer cancel()
	filters := bson.M{"_id": id}
	res := m.db.Collection("posts").FindOne(timeoutCtx, filters)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, posts.ErrPostNotFound
		}
		return nil, res.Err()
	}
	var doc PostDocument
	err := res.Decode(&doc)
	if err != nil {
		return nil, err
	}
	return m.deserialize(doc), nil
}

// Update modifies specific fields of a post identified by its ID.
func (m *MongoPostRepository) Update(ctx context.Context, postID string, update map[string]interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
	defer cancel()
	filter := bson.M{"_id": postID}
	updateInfo := bson.M{"$set": update}
	res, err := m.db.Collection("posts").UpdateOne(timeoutCtx, filter, updateInfo)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return posts.ErrPostNotFound
	}
	return nil
}

// serialize converts a posts.Post to a PostDocument for MongoDB storage.
func (m *MongoPostRepository) serialize(post posts.Post) *PostDocument {
	return &PostDocument{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Author:      post.Author,
		Slug:        post.Slug,
		Status:      post.Status,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}

// deserialize converts a PostDocument from MongoDB to a posts.Post.
func (m *MongoPostRepository) deserialize(doc PostDocument) *posts.Post {
	return &posts.Post{
		ID:          doc.ID,
		Title:       doc.Title,
		Content:     doc.Content,
		Author:      doc.Author,
		Slug:        doc.Slug,
		Status:      doc.Status,
		PublishedAt: doc.PublishedAt,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
}
