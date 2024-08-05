package ResponseHandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
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

func LoadMessages(serviceName string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("could not determine current file path")
	}

	// Get the directory containing the current file
	dir := filepath.Dir(filename)

	viper.SetConfigName("response")
	viper.SetConfigType("json")
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var data map[string][]Message

	if err := viper.Unmarshal(&data); err != nil {
		return err
	}
	serviceMessages, ok := data[serviceName]
	if !ok {
		return fmt.Errorf("service %s not found in configuration", serviceName)
	}

	response = make(map[int]Message)
	for _, msg := range serviceMessages {
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
		c.JSON(res.Status, res)
		return
	}
	c.JSON(500, gin.H{})
	return
}
