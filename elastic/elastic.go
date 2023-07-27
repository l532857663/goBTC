package elastic

import "fmt"

// 插入数据
func CreateData(_index, _id string, _source interface{}) error {
	filter := UrlFilter{
		Index:  _index,
		Id:     _id,
		Action: "_create",
	}
	res := &HitsInfo{}
	err := AskHttpJson(HttpPut, filter, _source, res)
	if res.Error != nil {
		return fmt.Errorf("CreateData error: [%s]%s", res.Error.Type, res.Error.Reason)
	}
	return err
}

// 修改数据
func UpdateData(_index, _id string, _source interface{}) error {
	filter := UrlFilter{
		Index:  _index,
		Id:     _id,
		Action: "_update",
	}
	res := &HitsInfo{}
	err := AskHttpJson(HttpPost, filter, _source, res)
	if res.Error != nil {
		return fmt.Errorf("UpdateData error: [%s]%s", res.Error.Type, res.Error.Reason)
	}
	return err
}

// 删除数据
