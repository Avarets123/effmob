package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func ParseStringToIntOrDefault(intStr string, defaultValue int) int {

	parsedInt, err := strconv.Atoi(intStr)

	if err != nil {
		return defaultValue
	}

	return parsedInt

}

func ConvertSliceToHashSet[T any](elems []T) map[any]bool {
	newMap := make(map[any]bool)
	for _, v := range elems {
		newMap[v] = true
	}
	return newMap

}

func MapMapToNeedStruct[T, M any](m M) T {

	var data T

	b, e := json.Marshal(m)

	fmt.Println(string(b))

	e = json.Unmarshal(b, &data)

	if e != nil {
		return data
	}

	return data

}
