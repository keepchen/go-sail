package vo

// Paginate paginator
type Paginate struct {
	Page     int `json:"page" query:"page" form:"page" path:"page" format:"number" example:"1"`                  //Page number, starting at 1
	PageSize int `json:"pageSize" query:"pageSize" form:"pageSize" path:"pageSize" format:"number" example:"20"` //Items per page
}
