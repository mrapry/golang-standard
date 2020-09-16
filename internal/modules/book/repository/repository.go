package repository

import (
	"master-service/internal/modules/book/repository/interfaces"
	"master-service/internal/modules/book/repository/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

// Repository parent
type Repository struct {
	readDB, writeDB *mongo.Database
	Book            interfaces.BookRepository
}

// NewRepository create new repository
func NewRepository(read, write *mongo.Database) *Repository {
	return &Repository{
		readDB: read, writeDB: write,
		Book: mongodb.NewBookRepo(read, write),
	}
}
