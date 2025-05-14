package domain

// PaginationParams represents the parameters used for paginating and sorting
// a collection of items in an API request.
//
// Fields:
// - Page: The current page number to retrieve.
// - PageSize: The number of items to include per page.
// - SortField: The field by which the items should be sorted.
// - SortOrder: The order of sorting; 1 for ascending and -1 for descending.
type PaginationParams struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"pageSize"`
	SortField string `form:"sortField"`
	SortOrder int    `form:"sortOrder"` // 1 for asc, -1 for desc
}
