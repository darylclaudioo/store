package db

import (
	"context"
	"hexagon-architecture/internal/domain/products/entity"
	"hexagon-architecture/internal/errors"
	"hexagon-architecture/internal/infrastructure"
	"hexagon-architecture/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Products struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name_product"`
	Stock int                `bson:"stock"`
}

func (db *DB) fromEntity(products entity.Products) Products {
	return Products{
		Name:  products.Name,
		Stock: products.Stock,
	}
}

func (products *Products) toEntity() *entity.Products {
	return &entity.Products{
		ID:    products.ID.Hex(),
		Name:  products.Name,
		Stock: products.Stock,
	}
}

func toEntities(p []Products) []*entity.Products {
	products := make([]*entity.Products, len(p))
	for i, product := range p {
		products[i] = product.toEntity()
	}
	return products
}

func (db *DB) GetProductByID(ctx context.Context, id string) (*entity.Products, int, error) {
	startTime := time.Now()

	ctx, span := infrastructure.Tracer().Start(ctx, "db:GetProductByID")
	defer span.End()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id":        _id,
		"deleted_at": bson.M{"$exists": false},
	}

	var pr Products
	err := db.db.Collection(db.products).FindOne(ctx, filter).Decode(&pr)
	if err == mongo.ErrNoDocuments {
		return nil, http.StatusNotFound, errors.ErrNotFoundProduct
	}

	endTime := time.Now()
	executionTime := endTime.Sub(startTime)

	log.Println(" Execution Time (Get Product By Id): %s\n", executionTime)
	return pr.toEntity(), http.StatusOK, nil
}

func (db *DB) GetCompanies(ctx context.Context, data entity.GetProductsRequest) ([]*entity.Products, *utils.Pagination, int, error) {
	startTime := time.Now()
	ctx, span := infrastructure.Tracer().Start(ctx, "db:GetCompanies")
	defer span.End()

	filter := make(map[string]interface{})

	if trimmedName := strings.TrimSpace(data.Name); trimmedName != "" {
		filter["name_product"] = bson.M{"$regex": primitive.Regex{Pattern: trimmedName, Options: "i"}}
	}

	limit := int64(data.Limit)
	skip := int64(data.Page*data.Limit - data.Limit)
	options := options.Find()

	options.SetLimit(limit)
	options.Skip = &skip

	cur, err := db.db.Collection(db.products).Find(ctx, filter, options)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, errors.ErrInternalDB
	}
	defer cur.Close(ctx)

	var products []Products
	for cur.Next(ctx) {
		var product Products
		err := cur.Decode(&product)
		if err != nil {
			return nil, nil, http.StatusInternalServerError, errors.ErrInternalDB
		}
		products = append(products, product)
	}
	if err := cur.Err(); err != nil {
		return nil, nil, http.StatusInternalServerError, errors.ErrInternalDB
	}

	total, err := db.db.Collection(db.products).CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, errors.ErrInternalDB
	}
	defer cur.Close(ctx)

	endTime := time.Now()
	executionTime := endTime.Sub(startTime)

	log.Println(" Execution Time (Get All Products): %s\n", executionTime)

	return toEntities(products), &utils.Pagination{
		Total:       int(total),
		Limit:       data.Limit,
		CurrentPage: data.Page,
		LastPage:    0,
	}, http.StatusOK, nil

}

func (db *DB) CreateProduct(ctx context.Context, product entity.Products) (*entity.Products, int, error) {
	startTime := time.Now()
	ctx, span := infrastructure.Tracer().Start(ctx, "db:CreateProduct")
	defer span.End()

	data := db.fromEntity(product)
	res, err := db.db.Collection(db.products).InsertOne(ctx, data)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.ErrInternalDB
	}

	data.ID = res.InsertedID.(primitive.ObjectID)

	endTime := time.Now()
	executionTime := endTime.Sub(startTime)

	log.Println(" Execution Time (Insert New Product): %s\n", executionTime)
	return data.toEntity(), http.StatusOK, nil
}

func (db *DB) UpdateProduct(ctx context.Context, id string, updateData entity.UpdateProductsRequest) (*entity.Products, int, error) {
	startTime := time.Now()
	ctx, span := infrastructure.Tracer().Start(ctx, "db:UpdateProduct")
	defer span.End()

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, http.StatusBadRequest, errors.ErrNotFoundProduct
	}

	filter := bson.M{
		"_id": _id,
	}

	update := bson.M{}
	if updateData.Name != "" {
		update["name_product"] = updateData.Name
	}
	if updateData.Stock != nil {
		update["stock"] = *updateData.Stock
	}

	updateQuery := bson.M{"$set": update}

	var pr Products
	err = db.db.Collection(db.products).FindOneAndUpdate(ctx, filter, updateQuery, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&pr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, http.StatusNotFound, errors.ErrNotFoundProduct
		}
		return nil, http.StatusInternalServerError, errors.ErrInternalDB
	}

	endTime := time.Now()
	executionTime := endTime.Sub(startTime)

	log.Println(" Execution Time (Update Product By Id): %s\n", executionTime)

	return pr.toEntity(), http.StatusOK, nil
}

func (db *DB) DeleteProduct(ctx context.Context, id string) (int, error) {
	startTime := time.Now()
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, errors.ErrNotFoundProduct
	}

	filter := bson.M{"_id": _id}

	result, err := db.db.Collection(db.products).DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return http.StatusNotFound, errors.ErrNotFoundProduct
		}
		return http.StatusInternalServerError, errors.ErrInternalDB
	}

	if result.DeletedCount == 0 {
		return http.StatusNotFound, errors.ErrNotFoundProduct
	}

	endTime := time.Now()
	executionTime := endTime.Sub(startTime)

	log.Println(" Execution Time (Delete Product By Id): %s\n", executionTime)
	return http.StatusOK, nil
}
