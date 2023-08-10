package utils

import (
	"encoding/json"
	"strconv"
)

func Map2Struct(theMap, theStruct interface{}) error {
	b, err := json.Marshal(theMap)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, theStruct)
	if err != nil {
		return err
	}
	return nil
}

func PageFilter(page, limit string) (int, int) {
	pageFrom, pageSize := 0, 20
	if page != "0" && page != "" && page != "1" {
		from, err := strconv.ParseUint(page, 10, 0)
		if err == nil {
			pageFrom = int(from)
		}
	}
	if limit != "0" && limit != "" {
		size, err := strconv.ParseUint(limit, 10, 0)
		if err == nil {
			pageSize = int(size)
		}
	}
	return pageFrom, pageSize
}
