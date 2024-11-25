package req

import (
	"encoding/json"
	"io"
)

func DecodeBody[T any](body io.ReadCloser) (T, error) {

	var resStruct T

	err := json.NewDecoder(body).Decode(&resStruct)

	return resStruct, err

}
