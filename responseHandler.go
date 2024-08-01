package ResponseHandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
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

func LoadMessages() error {
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

	var data struct {
		Handler []Message `json:"handler"`
	}

	if err := viper.Unmarshal(&data); err != nil {
		return err
	}
	response = make(map[int]Message)
	for _, msg := range data.Handler {
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
		Errors:  nil,
	}
}

func (r *Response) SendError(c *gin.Context, code int) {
	err := response[code]
	r.Message = err.Message
	r.Code = code
	c.JSON(err.Status, r)
	return
}

func (r *Response) SendSuccess(c *gin.Context) {
	r.Success = true
	r.Message = "OK"
	r.Code = 10000
	c.JSON(http.StatusOK, r)
	return
}
