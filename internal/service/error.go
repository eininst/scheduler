package service

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
)

const (
	ERROR_PARAMETER = 100
)

var ERROR_DATA_NOT_FOUND = NewServiceError("数据不存在")

type ServiceError struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

func (e *ServiceError) Error() string {
	jbytes, err := json.Marshal(*e)
	if err != nil {
		return err.Error()
	}

	return string(jbytes)
}

func NewServiceError(message interface{}, code ...int) error {
	cd := 0
	if len(code) > 0 {
		cd = code[0]
	}

	return &ServiceError{
		Code:    cd,
		Message: message,
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	codex := fiber.StatusInternalServerError
	contentType := fiber.MIMETextPlainCharsetUTF8
	var e *fiber.Error
	var sr *ServiceError

	if errors.As(err, &e) {
		codex = e.Code
	} else if errors.As(err, &sr) {
		codex = fiber.StatusBadRequest
		contentType = fiber.MIMEApplicationJSONCharsetUTF8
	}
	c.Set(fiber.HeaderContentType, contentType)

	return c.Status(codex).SendString(err.Error())
}
