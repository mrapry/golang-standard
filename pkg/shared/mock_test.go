package shared

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/stretchr/testify/assert"
)

func TestCreateHttpRequestBodyMock(t *testing.T) {
	type Test struct {
		Test string `json:"test"`
	}

	test := Test{}
	result := CreateHttpRequestBodyMock(test)

	assert.NotNil(t, result)
}

func TestSetMockerySharedResult(t *testing.T) {
	result := SetMockerySharedResult(golibshared.Result{Data: gofakeit.Word()})

	assert.NotNil(t, result)
}

func TestSetTestcaseName(t *testing.T) {
	name := SetTestcaseName(1, "positive set test case name")

	assert.Equal(t, SetTestcaseName(1, "positive set test case name"), name)
}
