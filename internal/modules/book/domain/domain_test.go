package domain

import (
	"testing"

	"github.com/Kamva/mgm/v3"
	"github.com/brianvoe/gofakeit"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBook_CollectionName(t *testing.T) {
	ID, _ := primitive.ObjectIDFromHex("5f637b7897623d6fa99923b0")
	type fields struct {
		IDField  mgm.IDField
		Name     string
		ISBN     string
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
				ISBN:     gofakeit.Name(),
				Version:  1,
				IsActive: true,
			},
			want: "book",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Book{
				IDField:  tt.fields.IDField,
				Name:     tt.fields.Name,
				ISBN:     tt.fields.ISBN,
				Version:  tt.fields.Version,
				IsActive: tt.fields.IsActive,
			}
			if got := m.CollectionName(); got != tt.want {
				t.Errorf("Book.CollectionName() = %v, want %v", got, tt.want)
			}
		})
	}
}
