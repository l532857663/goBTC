package oklink

import (
	"encoding/json"
	"fmt"
	"goBTC/global"

	"go.uber.org/zap"
)

func DecodeBodyDataOne[T BaseData](body []byte, bodyInfo *BaseResp[T]) (T, error) {
	funcName := "DecodeBodyDataOne"
	err := json.Unmarshal([]byte(body), &bodyInfo)
	if err != nil {
		global.LOG.Error("json Unmarshal error", zap.String("funcName", funcName), zap.Error(err))
		return nil, err
	}
	if bodyInfo.Code != "0" || len(bodyInfo.Data) == 0 {
		err := fmt.Errorf("Get response body error, code = %v, msg = '%v'", bodyInfo.Code, bodyInfo.Msg)
		global.LOG.Error("Get response body error error", zap.String("funcName", funcName), zap.Error(err))
		return nil, err
	}
	data := bodyInfo.Data[0]
	global.LOG.Info("DecodeBody info success", zap.String("funcName", funcName))
	global.LOG.Debug("DecodeBody info:", zap.Any("address info", data))
	return data, nil
}

func DecodeBodyDataAll[T BaseData](body []byte, bodyInfo *BaseResp[T]) ([]T, error) {
	funcName := "DecodeBodyDataAll"
	err := json.Unmarshal([]byte(body), &bodyInfo)
	if err != nil {
		global.LOG.Error("json Unmarshal error", zap.String("funcName", funcName), zap.Error(err))
		return nil, err
	}
	if bodyInfo.Code != "0" || len(bodyInfo.Data) == 0 {
		err := fmt.Errorf("Get response body error, code = %v, msg = '%v'", bodyInfo.Code, bodyInfo.Msg)
		global.LOG.Error("Get response body error error", zap.String("funcName", funcName), zap.Error(err))
		return nil, err
	}
	data := bodyInfo.Data
	global.LOG.Info("DecodeBody info success", zap.String("funcName", funcName))
	global.LOG.Debug("DecodeBody info:", zap.Any("address info", data))
	return data, nil
}
