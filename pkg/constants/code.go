package constants

//ICodeType 错误码类型接口
type ICodeType interface {
	String() string
	Int() int
}

//CodeType 错误码类型
type CodeType int
