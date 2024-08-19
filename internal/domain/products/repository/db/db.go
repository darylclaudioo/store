package db

import "go.mongodb.org/mongo-driver/mongo"

type DB struct {
	db       *mongo.Database
	products string
}

func New(db *mongo.Database, products string) *DB {
	return &DB{
		db:       db,
		products: products,
	}
}
