package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"stock-api/infrastructure/core/domain"
	"stock-api/infrastructure/core/port"
	"stock-api/infrastructure/response"
)

type StockHandler struct {
	stockService           port.StockService
	serviceBestInvestments port.BestInvestmentsService
	workerPool             chan struct{}
}

func NewStockHandler(service port.StockService, service_best_investments port.BestInvestmentsService, maxWorkers int) *StockHandler {
	return &StockHandler{stockService: service, serviceBestInvestments: service_best_investments, workerPool: make(chan struct{}, maxWorkers)}
}

// FindStocks handles the HTTP request to retrieve a list of stocks.
// It supports pagination, sorting, and filtering.
//
// @Summary Retrieve stocks
// @Description Retrieves a list of stocks based on pagination, sorting, and optional filters.
// @Tags stocks
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param size query int false "Page size for pagination"
// @Param sort query string false "Sorting criteria (e.g., 'name asc')"
// @Param filters body domain.Filters false "Filters to apply to the stock search"
// @Success 200 {object} []domain.Stock "List of stocks"
// @Failure 400 {object} response.ErrorResponse "Invalid parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve stocks"
// @Router /stocks [get]
func (h *StockHandler) FindStocks(c *gin.Context) {
	// Retrieves the pagination parameters from the query string
	// and binds them to the PaginationParams struct.
	// The query parameters are expected to be in the format:
	// ?page=1&size=10&sort=name asc
	var pagination domain.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, "Invalid parameters")
		return
	}

	// Retrieves the filters from the request body and binds them to the Filters struct.
	// The filters are expected to be in JSON format.
	// If no filters are provided, an empty Filters struct is created.
	// The filters are used to filter the stocks based on specific criteria.

	var requestBody domain.FilterRequest

	// Bind the JSON from the request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.BadRequest(c, "Invalid filters")
		return
	}

	filters := requestBody.Filters
	if filters == nil {
		filters = make(domain.Filters) // Initialize if no filters are provided
	}

	// Calls the service to find stocks based on the pagination and filters.
	stocks, total, err := AsyncManyOperation(c, h.workerPool, func() ([]domain.Stock, int, error) {
		return h.stockService.Find(c.Request.Context(), pagination, filters)
	})

	if err != nil {
		response.InternalServerError(c, "Failed to retrieve stocks")
		return
	}

	resp := response.ToStockResponse(stocks, pagination.PageSize, total, pagination.SortField)

	// Returns the list of stocks in the response with a 200 status code.
	response.Success(c, 200, resp)
}

// GetStockRecommendations handles the HTTP request to retrieve stock recommendations.
// It uses a default limit of 5 recommendations unless specified in the query parameters.
//
// Query Parameters:
// - limit: (optional) The maximum number of recommendations to return.
//
// Responses:
// - 200: Returns a JSON response with the list of stock recommendations.
// - 500: Returns an internal server error if there is an issue retrieving the stocks.
func (h *StockHandler) GetStockRecommendations(c *gin.Context) {
	limit := 5
	if c.Query("limit") != "" {
		limit, _ = strconv.Atoi(c.Query("limit"))
	}

	pagination := domain.PaginationParams{
		Page:     1,
		PageSize: 5000,
	}

	filters := domain.Filters{
		"classifications": domain.Filter{
			Value: []string{
				"High-Risk Speculative",
				"Bearish Signal",
				"Analyst Negative",
			},
			MatchMode: "notOverlap",
		},
	}

	// Calls the service to find stocks based on the pagination and filters.
	stocks, _, err := AsyncManyOperation(c, h.workerPool, func() ([]domain.Stock, int, error) {
		return h.stockService.Find(c.Request.Context(), pagination, filters)
	})

	if err != nil {
		response.InternalServerError(c, "Failed to retrieve stocks")
		return
	}

	recommendations := h.serviceBestInvestments.GetStockRecommendations(stocks, limit)

	response.Success(c, 200, recommendations)
}
