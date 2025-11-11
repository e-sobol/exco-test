package fetcher

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ReadBody(ctx *gin.Context, body any) error {
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(body)
	if err != nil {
		return err
	}
	err = validator.New().Struct(body)
	if err != nil {
		return err
	}
	return nil
}

func UnwrapPointerOrDefault[T any](val *T, defaultVal T) T {
	if val != nil {
		return *val
	}
	return defaultVal
}
