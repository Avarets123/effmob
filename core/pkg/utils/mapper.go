package utils

func MapSlice[T, RT any](data []T, mapper func(T) RT) []RT {
	var resData []RT
	for _, v := range data {
		resData = append(resData, mapper(v))
	}
	return resData

}

func ReduceSlice[T, RT any](data []T, reducer func(acc RT, el T, ind int) RT) RT {
	var resData RT

	for i, v := range data {
		resData = reducer(resData, v, i)
	}

	return resData

}
