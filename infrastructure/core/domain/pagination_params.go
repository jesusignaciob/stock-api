package domain

type PaginationParams struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"pageSize"`
	SortField string `form:"sortField"`
	SortOrder int    `form:"sortOrder"` // 1 for asc, -1 for desc
}
