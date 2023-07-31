package ord

import (
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
	res, err := elastic.GetDataByFilter(elastic.InscribeInfoType, searchInfo)
	if err != nil {
		return nil, err
	}
	if res.Hits.Total.Value == 0 {
		return nil, nil
	}
	return &res.Hits, nil
}
