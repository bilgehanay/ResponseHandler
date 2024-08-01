package ResponseHandler

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

//go:embed response.json
var responseFile []byte

type Message struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  int    `json:"status"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TraceId string      `json:"traceId,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
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

var errorMap map[int]Message
var successMessage Message

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

func (r Response) SendError(c *gin.Context, code int) {
	err := errorMap[code]
	r.Message = err.Message
	r.Code = code
	fmt.Println(errorMap)
	c.JSON(err.Status, r)
	return
}

func (r Response) SendSuccess(c *gin.Context) {
	r.Message = "OK"
	r.Code = 10000
	c.JSON(http.StatusOK, r)
	return
}

func LoadMessages() error {

	viper.SetConfigType("json")

	if err := viper.ReadConfig(bytes.NewReader(responseFile)); err != nil {
		return fmt.Errorf("could not open message file: %v", err)
	}

	var data struct {
		Errors  []Message `json:"errors"`
		Success Message   `json:"success"`
	}

	if err := viper.Unmarshal(&data); err != nil {
		return fmt.Errorf("could not decode message file: %v", err)
	}

	errorMap = make(map[int]Message)
	for _, msg := range data.Errors {
		errorMap[msg.Code] = msg
	}

	successMessage = data.Success
	return nil
}

func HandleError(c *gin.Context, code int, traceId string, data interface{}, errs interface{}) {
	if msg, ok := errorMap[code]; ok {
		response := ErrorResponse{
			Success: false,
			Code:    msg.Code,
			Message: msg.Message,
			TraceId: traceId,
			Data:    data,
			Errors:  errs,
		}
		c.JSON(msg.Status, response)
	} else {
		c.JSON(500, gin.H{"success": false, "message": "Internal Server Error"})
	}
}

func HandleSuccess(c *gin.Context, data interface{}, count int) {
	c.JSON(successMessage.Status, gin.H{
		"success": true,
		"errors":  nil,
		"code":    successMessage.Code,
		"message": successMessage.Message,
		"data":    data,
		"count":   count,
	})
}
