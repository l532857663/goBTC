package resp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	TimeStamp string      `json:"timeStamp"`
	Sign      string      `json:"sign"`
}

const (
	SUCCESS = 200
	WARN    = 400
	ERROR   = 500
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 构造待签名结构，只对数据和时间戳签名
	var r Response

	// 增加返回结构中其他数据信息
	r.Status = strconv.Itoa(code)
	r.Message = msg
	r.Data = data
	r.TimeStamp = strconv.FormatInt(time.Now().Unix(), 10)

	c.JSON(http.StatusOK, r)
}

func Ok(c *gin.Context) {
	Result(SUCCESS, nil, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, nil, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

func OkWithCodeData(code int, data interface{}, c *gin.Context) {
	Result(code, data, "操作成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, nil, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, nil, message, c)
}

func FailWithCodeMessage(code int, message string, c *gin.Context) {
	Result(code, nil, message, c)
}
func FailParamError(c *gin.Context) {
	Result(ERROR, nil, "required param error", c)
}
func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func Custom(response interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, response)
}

func UnmarshalResponse(body string) (*Response, error) {
	result := &Response{}
	err := json.Unmarshal([]byte(body), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
