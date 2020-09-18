package main

//DomainTemplate for default generator domain
//{{$.LibraryAddress}} = address library
//{{$.module}} = name of module (lowercase)
//{{$.Petik}} = ` (symbol petik)
//{{clean (upper $.module)}} = name of module (capital)
const DomainTemplate = `package domain

import (
	"github.com/Kamva/mgm/v3"
	"{{$.LibraryAddress}}/golibshared"
)

const (
	//Collection{{clean (upper $.module)}} name constanta for domain module {{$.module}}
	Collection{{clean (upper $.module)}} = "{{$.module}}"
)

// {{clean (upper $.module)}} structure
type {{clean (upper $.module)}} struct {
	mgm.IDField {{$.Petik}}bson:",inline"{{$.Petik}}
	Name        string {{$.Petik}}bson:"name" json:"name"{{$.Petik}}
	Version     int    {{$.Petik}}bson:"version" json:"version,omitempty"{{$.Petik}}
	IsActive    bool   {{$.Petik}}bson:"isActive" json:"isActive,omitempty"{{$.Petik}}
}

// CollectionName for {{$.module}} model
func (m *{{clean (upper $.module)}}) CollectionName() string {
	return Collection{{clean (upper $.module)}}
}

// Filter model
type Filter struct {
	golibshared.Filter
}

// FieldMap mapping json to column name
var FieldMap = map[string]string{
	"id":   "_id",
	"name": "name",
}
`

//DomainTestTemplate for default generator domain test
//{{$.module}} = name of module (lowercase)
//{{clean (upper $.module)}} = name of module (capital)
const DomainTestTemplate = `package domain

import (
	"testing"

	"github.com/Kamva/mgm/v3"
	"github.com/brianvoe/gofakeit"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test{{clean (upper $.module)}}_CollectionName(t *testing.T) {
	ID, _ := primitive.ObjectIDFromHex("5f637b7897623d6fa99923b0")
	type fields struct {
		IDField  mgm.IDField
		Name     string
		Version  int
		IsActive bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test case get name of collection in constant",
			fields: fields{
				IDField:  mgm.IDField{ID: ID},
				Name:     gofakeit.Name(),
				Version:  1,
				IsActive: true,
			},
			want: "{{$.module}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &{{clean (upper $.module)}}{
				IDField:  tt.fields.IDField,
				Name:     tt.fields.Name,
				Version:  tt.fields.Version,
				IsActive: tt.fields.IsActive,
			}
			if got := m.CollectionName(); got != tt.want {
				t.Errorf("{{clean (upper $.module)}}.CollectionName() = %v, want %v", got, tt.want)
			}
		})
	}
}
`
