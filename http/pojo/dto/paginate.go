package dto

// Paginate paginator
type Paginate struct {
	Page       int   `json:"page" format:"number" example:"1" validate:"required"`        //Page number, starting at 1
	PageSize   int   `json:"pageSize" format:"number" example:"10" validate:"required"`   //Items per page
	TotalCount int64 `json:"totalCount" format:"number" example:"20" validate:"required"` //Total count
}
