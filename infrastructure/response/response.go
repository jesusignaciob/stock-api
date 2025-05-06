package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(ctx *gin.Context, status int, data interface{}) {
	ctx.IndentedJSON(status, JsonResponse{
		Success: true,
		Data:    data,
	})
}

func Error(ctx *gin.Context, status int, err string) {
	ctx.IndentedJSON(status, JsonResponse{
		Success: false,
		Error:   err,
	})
}

// Funciones adicionales útiles
func Created(ctx *gin.Context, data interface{}) {
	Success(ctx, http.StatusCreated, data)
}

func BadRequest(ctx *gin.Context, err string) {
	Error(ctx, http.StatusBadRequest, err)
}

func NotFound(ctx *gin.Context, err string) {
	Error(ctx, http.StatusNotFound, err)
}

func InternalServerError(ctx *gin.Context, err string) {
	Error(ctx, http.StatusInternalServerError, err)
}
