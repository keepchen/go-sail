package vo

// Paginate 分页器
type Paginate struct {
	Page     int `json:"page" query:"page" form:"page" path:"page" format:"number" example:"1"`                  //页码，从1开始
	PageSize int `json:"pageSize" query:"pageSize" form:"pageSize" path:"pageSize" format:"number" example:"20"` //页面条数
}
