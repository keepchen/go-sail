package api

import (
	"fmt"
	"time"

	"github.com/keepchen/go-sail/v3/constants"
)

var (
	anotherErrNoneCode constants.ICodeType = constants.ErrNone         //被改写后的成功code码
	emptyDataField     interface{}         = nil                       //空data字段
	forceHttpCode200                       = false                     //强制使用200作为http的状态码
	timezone                               = constants.DefaultTimeZone //时区
)

var (
	loc *time.Location
)

// Option 配置项
type Option struct {
	//成功code码，默认成功code码为0，配置此项后，成功code码将使用这个值。
	ErrNoneCode constants.ICodeType
	//成功提示语，默认成功提示语为SUCCESS，配置此项后，成功提示语将使用这个值。
	ErrNoneCodeMsg string
	//空data序列化结构，默认返回的data字段为空时为null值，配置此项后，空data序列化格式将使用这个值。
	EmptyDataStruct int
	//强制使用200作为http的状态码，配置此项后，http状态码将不从业务code码中解析。
	//
	//注意，调用Status()方法和SendWithCode()方法时的优先级高于此项配置。
	ForceHttpCode200 bool
	//时区
	Timezone string
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
	if opt.ForceHttpCode200 {
		forceHttpCode200 = opt.ForceHttpCode200
	}
	if len(opt.Timezone) > 0 {
		timezone = opt.Timezone
	}

	lc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(fmt.Errorf("[GO-SAIL] can not load location: %s", timezone))
	}
	loc = lc
}
