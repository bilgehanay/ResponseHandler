package ResponseHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Message struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  int    `json:"status"`
}

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TraceId string      `json:"traceId,omitempty"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

var response map[int]Message

func LoadMessages(responsePath string) error {

	viper.SetConfigFile(responsePath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var data struct {
		Response []Message `json:"response"`
	}

	if err := viper.Unmarshal(&data); err != nil {
		return err
	}

	response = make(map[int]Message)
	for _, msg := range data.Response {
		response[msg.Code] = msg
	}

	return nil
}

func New() *Response {
	return &Response{
		Success: false,
		Code:    0,
		Message: "",
		TraceId: "",
		Data:    nil,
		Count:   0,
		Errors:  nil,
	}
}

func (r *Response) SendResponse(c *gin.Context, code int) {
	if res, ok := response[code]; ok {
		if res.Status < 300 {
			r.Success = true
		}
		r.Message = res.Message
		r.Code = res.Code
		c.JSON(res.Status, r)
		return
	}
	c.JSON(500, gin.H{})
	return
}
