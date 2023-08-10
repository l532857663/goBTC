package ord

import (
	"fmt"
	"goBTC/elastic"
)

func GetInscribeIsExist(txId string) (string, error) {
	searchInfo := elastic.SearchInfo{
		Query: &elastic.Query{},
	}
	searchInfo.Query.Match = make(map[string]interface{})
	searchInfo.Query.Match["owner_output"] = txId

	res, err := elastic.GetDataByFilter(elastic.InscribeInfoType, searchInfo)
	if err != nil {
		return "", err
	}
	if res.Hits.Total.Value == 0 || len(res.Hits.Hits) == 0 {
		return "", nil
	}
	return res.Hits.Hits[0].Id, nil
}

func GetUnSyncOrdToken() (*elastic.Hits, error) {
	searchInfo := elastic.SearchInfo{
		Query: &elastic.Query{},
	}
	searchInfo.Query.Match = make(map[string]interface{})
	searchInfo.Query.Match["sync_state"] = elastic.StateSyncIsFalse
	res, err := elastic.GetDataByFilter(elastic.OrdTokenType, searchInfo)
	if err != nil {
		return nil, err
	}
	if res.Hits.Total.Value == 0 {
		return nil, nil
	}
	return &res.Hits, nil
}

func GetBalanceInfo(esIndex string, pageFrom, pageSize int, filter map[string]interface{}) ([]interface{}, error) {
	fmt.Printf("wch---- filter: %+v\n", filter)
	searchInfo := elastic.SearchInfo{
		From: pageFrom,
		Size: pageSize,
		Query: &elastic.Query{
			Match: filter,
		},
	}

	res, err := elastic.GetDataByFilter(esIndex, searchInfo)
	if err != nil {
		return nil, err
	}
	if res.Hits.Total.Value == 0 || len(res.Hits.Hits) == 0 {
		return nil, nil
	}
	fmt.Printf("wch----- res: %+v\n", res.Hits.Total.Value)
	var list []interface{}
	for _, info := range res.Hits.Hits {
		list = append(list, info.Source)
	}
	return list, nil
}
