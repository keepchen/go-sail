package dto

// IResponse common response interface
type IResponse interface {
	//GetData retrieves the actual data field value
	GetData() any
}
