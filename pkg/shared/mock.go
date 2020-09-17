package shared

import (
	"encoding/json"
	"fmt"

	"github.com/mrapry/go-lib/golibshared"
)

//CreateHttpRequestBodyMock create http body mock
func CreateHttpRequestBodyMock(structure interface{}) string {
	json, _ := json.Marshal(structure)
	result := string(json)

	return result
}

//SetMockerySharedResult set shared result mock
func SetMockerySharedResult(result interface{}) <-chan golibshared.Result {
	sharedResult := result.(golibshared.Result)

	// simulasiin untuk set channel shared result
	resultShared := func() <-chan golibshared.Result {
		output := make(chan golibshared.Result)
		go func() { output <- sharedResult }()
		return output
	}()

	return resultShared
}

//SetTestcaseName set testcase name to prevent tech debt
func SetTestcaseName(number int, description string) string {
	return fmt.Sprintf("Testcase #%v : %s", number, description)
}
