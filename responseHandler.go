package ResponseHandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

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

var errorMap map[int]Message
var successMessage Message

func LoadMessages() error {

	executablePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("could not get executable path: %v", err)
	}
	filePath := filepath.Join(executablePath, "response.json")
	viper.SetConfigFile(filePath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
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
