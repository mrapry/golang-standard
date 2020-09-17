package mongodb

import (
	"context"
	"master-service/internal/modules/book/domain"
	"master-service/internal/modules/book/repository/interfaces"

	db "github.com/Kamva/mgm/v3"
	"github.com/Kamva/mgm/v3/operator"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/mrapry/go-lib/tracer"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type bookRepoMongDB struct {
	readDB, writeDB *mongo.Database
}

// NewBookRepo create new book repository
func NewBookRepo(readDB, writeDB *mongo.Database) interfaces.BookRepository {
	return &bookRepoMongDB{readDB, writeDB}
}

func (r *bookRepoMongDB) FindAll(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.Book{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.readDB, collName)

		// set offset
		filter.CalculateOffset()

		// set sort
		filter.SetSort()

		// set order by
		orderBy := filter.SetOrderBy(domain.FieldMap)

		// set default query
		where := []bson.M{}

		// set search
		fields := []string{"name"}
		where = filter.SetSearch(where, fields)

		// set show all
		if !filter.ShowAll {
			where = append(where, bson.M{"isActive": true})
		}

		// set option
		limit := cast.ToInt64(filter.Limit)
		offset := cast.ToInt64(filter.Offset)
		findOptions := &options.FindOptions{
			Limit: &limit,
			Skip:  &offset,
			Sort:  orderBy,
		}

		// set query
		query := bson.M{operator.And: where}

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.Find,
			Filter:     query,
			Sort:       findOptions.Sort,
			Limit:      *findOptions.Limit,
			Skip:       *findOptions.Skip,
		}
		trace.SetTags(ctx)

		var book = []*domain.Book{}
		if err := coll.SimpleFind(&book, query, findOptions); err != nil {
			tracer.SetError(ctx, err)
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: book}
	}()

	return output
}

func (r *bookRepoMongDB) Count(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.Book{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.readDB, collName)

		// set default query
		where := []bson.M{}

		// set search
		fields := []string{"name"}
		where = filter.SetSearch(where, fields)

		// set show all
		if !filter.ShowAll {
			where = append(where, bson.M{"isActive": true})
		}

		// set query
		query := bson.M{operator.And: where}

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.CountDocument,
			Filter:     query,
		}
		trace.SetTags(ctx)

		count, err := coll.CountDocuments(ctx, query)
		if err != nil {
			tracer.SetError(ctx, err)
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: count}
	}()

	return output
}

func (r *bookRepoMongDB) Find(ctx context.Context, obj domain.Book) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.Book{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.readDB, collName)

		// set data to bson M
		query := golibshared.ToBSON(obj)

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.FindOne,
			Filter:     query,
		}
		trace.SetTags(ctx)

		if err := coll.First(query, model); err != nil {
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: model}
	}()

	return output
}

func (r *bookRepoMongDB) FindByID(ctx context.Context, ID string) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.Book{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.readDB, collName)

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.FindOne,
			Filter:     ID,
		}
		trace.SetTags(ctx)

		if err := coll.FindByID(ID, model); err != nil {
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: model}
	}()

	return output
}

func (r *bookRepoMongDB) Save(ctx context.Context, data *domain.Book) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.Book{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.writeDB, collName)

		// set version
		data.Version = data.Version + 1

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.UpdateOne,
			Filter:     data,
		}
		trace.SetTags(ctx)

		if err := coll.Update(data); err != nil {
			tracer.SetError(ctx, err)
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: data}
	}()

	return output
}

func (r *bookRepoMongDB) Insert(ctx context.Context, newData *domain.Book) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.Book{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.writeDB, collName)

		// set version
		newData.Version = newData.Version + 1

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.InsertOne,
			Filter:     newData,
		}
		trace.SetTags(ctx)

		if err := coll.Create(newData); err != nil {
			tracer.SetError(ctx, err)
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: newData}

	}()

	return output
}
