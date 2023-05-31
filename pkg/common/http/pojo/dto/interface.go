package dto

// IResponse 统一返回接口
type IResponse interface {
	//GetData 获取data字段实际数据
	GetData() interface{}
}
