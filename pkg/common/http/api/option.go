package api

import "github.com/keepchen/go-sail/v2/pkg/constants"

var (
	anotherErrNoneCode constants.ICodeType = constants.ErrNone //被改写后的成功code码
	emptyDataField     interface{}         = nil               //空data字段
)

// Option 配置项
type Option struct {
	ErrNoneCode     constants.ICodeType //成功code码
	ErrNoneCodeMsg  string              //成功提示语
	EmptyDataStruct int                 //空data序列化结构
}

const (
	DefaultEmptyDataStructNull   int = 0 //空(null)
	DefaultEmptyDataStructObject int = 1 //大(花)括号
	DefaultEmptyDataStructArray  int = 2 //中(方)括号
	DefaultEmptyDataStructString int = 3 //空字符串
)

// SetupOption 设置选项
//
// 目前支持设定:
//
// 1.成功code码及对应code信息
//
// 2.空数据序列化结构
func SetupOption(opt Option) {
	if opt.ErrNoneCode != nil {
		constants.RegisterCode(constants.CodeType(opt.ErrNoneCode.Int()), opt.ErrNoneCodeMsg)
		anotherErrNoneCode = opt.ErrNoneCode
	}
	switch opt.EmptyDataStruct {
	case DefaultEmptyDataStructNull:
		emptyDataField = nil
	case DefaultEmptyDataStructObject:
		emptyDataField = struct{}{}
	case DefaultEmptyDataStructArray:
		emptyDataField = []bool{}
	case DefaultEmptyDataStructString:
		emptyDataField = ""
	default:
		emptyDataField = nil
	}
}
