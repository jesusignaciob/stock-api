package domain

type Filter struct {
	Value     interface{} `json:"value"`
	MatchMode string      `json:"matchMode"`
}

type Filters map[string]Filter

// Estructura para el body completo
type FilterRequest struct {
	Filters Filters `json:"filters"`
}
