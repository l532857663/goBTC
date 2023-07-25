package ord

import (
	"goBTC/elastic"
)

func GetInscribeIsExist(txId string) (bool, error) {
	searchInfo := elastic.SearchInfo{
		Source: []string{"owner_output"},
		Query:  &elastic.Query{},
	}
	searchInfo.Query.Match = make(map[string]interface{})
	searchInfo.Query.Match["owner_output"] = txId

	res, err := elastic.GetDataByFilter(elastic.InscribeInfoType, searchInfo)
	if err != nil {
		return false, err
	}
	if res.Hits.Total.Value == 0 {
		return false, nil
	}
	return true, nil
}
