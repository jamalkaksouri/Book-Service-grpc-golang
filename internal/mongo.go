package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const (
	bookCollection = "book"
)

type MongoBookRepository struct {
	counter    BookId
	mu         sync.Mutex
	collection *mongo.Collection
}

func NewMongoBookRepository(db *mongo.Database) *MongoBookRepository {
	return &MongoBookRepository{
		collection: db.Collection(bookCollection),
	}
}

func (r *MongoBookRepository) CreateBook(ctx context.Context, book *Book) (BookId, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++

	_, err := r.collection.InsertOne(ctx, book)
	if err != nil {
		return 0, err
	}
	return r.counter, err
}

func (r *MongoBookRepository) RetrieveBook(ctx context.Context, id BookId) (*Book, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var book Book
	err := r.collection.FindOne(ctx, map[string]BookId{"id": id}).Decode(&book)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *MongoBookRepository) UpdateBook(ctx context.Context, book *Book) error {
	_, err := r.collection.UpdateOne(ctx, map[string]BookId{"id": book.Id}, map[string]interface{}{"$set": book})
	return err
}

func (r *MongoBookRepository) DeleteBook(ctx context.Context, id BookId) error {
	_, err := r.collection.DeleteOne(ctx, map[string]BookId{"id": id})
	return err
}

func (r *MongoBookRepository) ListBook(ctx context.Context, offset int64, limit int64) ([]*Book, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var books []*Book
	cursor, err := r.collection.Find(ctx, map[string]interface{}{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)
	for cursor.Next(ctx) {
		var book Book
		if err := cursor.Decode(&book); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}
	return books, nil
}
