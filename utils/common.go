package utils

import "encoding/json"

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
