package req

import (
	"effect-mobile/pkg/res"
	"io"
)

func HandleBody[T any](body io.ReadCloser) (T, res.Error) {

	data, err := DecodeBody[T](body)

	if err != nil {
		return data, res.NewError(err, 422)
	}

	err = IsValidaStruct(data)

	if err != nil {
		return data, res.NewError(err, 422)

	}

	return data, nil

}
