package response

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationError(err error) gin.H {
	errors := make(map[string]string)

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return gin.H{"error": "invalid request"}
	}

	for _, fe := range ve {
		field := toSnake(fe.Field())

		switch fe.Tag() {
		case "required":
			errors[field] = "is required"
		case "min":
			errors[field] = "must be at least " + fe.Param() + " characters"
		case "email":
			errors[field] = "must be a valid email"
		default:
			errors[field] = "is invalid"
		}
	}

	return gin.H{"errors": errors}
}

func toSnake(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
