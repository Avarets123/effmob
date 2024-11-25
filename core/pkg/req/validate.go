package req

import "github.com/go-playground/validator/v10"

func IsValidaStruct[T any](data T) error {
	return validator.New().Struct(data)
}
