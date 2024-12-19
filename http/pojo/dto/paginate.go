package dto

// Paginate 分页器
type Paginate struct {
	Page       int   `json:"page" format:"number" example:"1" validate:"required"`        //页码，从1开始
	PageSize   int   `json:"pageSize" format:"number" example:"10" validate:"required"`   //页面条数
	TotalCount int64 `json:"totalCount" format:"number" example:"20" validate:"required"` //总条数
}
