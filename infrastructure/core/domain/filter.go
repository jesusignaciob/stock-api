package domain

// Filter represents a single filter criterion with a value and a match mode.
// The Value field holds the value to filter by, and the MatchMode field specifies
// the type of matching to apply (e.g., exact, contains, etc.).
//
// Example JSON representation:
//
//	{
//	  "value": "example",
//	  "matchMode": "contains"
//	}
type Filter struct {
	Value     interface{} `json:"value"`
	MatchMode string      `json:"matchMode"`
}

// Filters is a map where each key represents a field name, and the value is a Filter
// that defines the filtering criteria for that field.
//
// Example JSON representation:
//
//	{
//	  "fieldName": {
//	    "value": "example",
//	    "matchMode": "contains"
//	  }
//	}
type Filters map[string]Filter

// FilterRequest represents the structure of a request body containing multiple filters.
// The Filters field is a map of field names to their respective filtering criteria.
//
// Example JSON representation:
//
//	{
//	  "filters": {
//	    "fieldName": {
//	      "value": "example",
//	      "matchMode": "contains"
//	    }
//	  }
//	}
type FilterRequest struct {
	Filters Filters `json:"filters"`
}
