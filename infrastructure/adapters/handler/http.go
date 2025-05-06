package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"stock-api/infrastructure/core/port"
	"stock-api/infrastructure/response"
)

type StockHandler struct {
	service port.StockService
}

func NewStockHandler(service port.StockService) *StockHandler {
	return &StockHandler{service: service}
}

func (h *StockHandler) FindAllStocks(c *gin.Context) {
	order := c.Query("order")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.BadRequest(c, "Invalid page parameter")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter")
		return
	}

	stocks, err := h.service.FindAllStocks(c, order, page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve stocks")
		return
	}

	response.Success(c, 200, stocks)
}
