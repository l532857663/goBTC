package ord

import (
	"goBTC/elastic"
	"goBTC/utils"
)

func GetInscribeIsExist(txId, output string) (string, error) {
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
	// 判断outout内容
	for _, hit := range res.Hits.Hits {
		info := &elastic.InscribeInfo{}
		err = utils.Map2Struct(hit.Source, info)
		if err != nil {
			return "", nil
		}
		if info.OwnerOutput == output {
			return hit.Id, nil
		}
	}
	return "", nil
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
	var list []interface{}
	for _, info := range res.Hits.Hits {
		list = append(list, info.Source)
	}
	return list, nil
}

func GetUnSyncOrdInscribe() (*elastic.Hits, error) {
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
