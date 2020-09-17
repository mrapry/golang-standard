package domain

import (
	"github.com/Kamva/mgm/v3"
	"github.com/mrapry/go-lib/golibshared"
)

const (
	//CollectionBook name constanta for domain module book
	CollectionBook = "book"
)

// Book structure
type Book struct {
	mgm.IDField `bson:",inline"`
	Name        string `bson:"name" json:"name"`
	ISBN        string `bson:"isbn" json:"isbn,omitempty"`
	Version     int    `bson:"version" json:"version,omitempty"`
	IsActive    bool   `bson:"isActive" json:"isActive,omitempty"`
}

// CollectionName for book model
func (m *Book) CollectionName() string {
	return CollectionBook
}

// Filter model
type Filter struct {
	golibshared.Filter
}

// FieldMap mapping json to column name
var FieldMap = map[string]string{
	"id":   "_id",
	"name": "name",
	"isbn": "isbn",
}
