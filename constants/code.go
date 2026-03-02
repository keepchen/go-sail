package constants

import "sync"

// ICodeType 错误码类型接口
type ICodeType interface {
	String(language ...string) string
	Int() int
}

// CodeType 错误码类型
type CodeType int

// MMBox 语言码-错误码映射类型
type MMBox map[LanguageCode]map[ICodeType]string

// errorCodeTypeMsgMap 使用 sync.Map 实现并发安全
// 外层：LanguageCode -> *sync.Map
// 内层：ICodeType -> string
type errorCodeTypeMsgMap struct {
	maps sync.Map // map[LanguageCode]*sync.Map
}

var (
	ctm  *errorCodeTypeMsgMap
	once sync.Once
)

// RegisterCodeTable 注册错误码表
//
// i18nMsg key为错误码值，value为错误信息
//
// 与 RegisterCodeSingle 不同的是，此方法是覆盖对应语言的整个code表。
//
// # 此方法适用于【固定】注入错误码表或需要覆盖默认错误码表的场景
func RegisterCodeTable(language LanguageCode, i18nMsg map[ICodeType]string) {
	m := &sync.Map{}
	for k, v := range i18nMsg {
		m.Store(k, v)
	}
	ctm.maps.Store(language, m)
}

// RegisterCodeSingle 注册单个错误码
//
// 与 RegisterCodeTable 不同的是，此方法更细粒度的注入错误码及对应的错误信息而不是整个覆盖。
//
// 当code码重复时，后者覆盖前者
//
// # 此方法适用于【动态】注入单个错误码的场景
func RegisterCodeSingle(language LanguageCode, code ICodeType, msg string) {
	// 确保内层存在
	val, _ := ctm.maps.LoadOrStore(language, &sync.Map{})
	inner := val.(*sync.Map)
	inner.Store(code, msg)
}

// GetCodeMsg 获取错误码对应的错误信息
func GetCodeMsg(language LanguageCode, code ICodeType) (string, bool) {
	if val, ok := ctm.maps.Load(language); ok {
		if inner, ok2 := val.(*sync.Map); ok2 {
			if msg, ok3 := inner.Load(code); ok3 {
				return msg.(string), true
			}
		}
	}

	return "", false
}

func init() {
	once.Do(func() {
		ctm = &errorCodeTypeMsgMap{}

		for language, codeMsg := range initErrorCodeMsgMap {
			RegisterCodeTable(language, codeMsg)
		}
	})
}
