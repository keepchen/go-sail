package constants

import "sync"

// ICodeType 错误码类型接口
type ICodeType interface {
	String(language ...string) string
	Int() int
}

// CodeType 错误码类型
type CodeType int

type errorCodeTypeMsgMap struct {
	mux  *sync.RWMutex
	maps map[LanguageCode]map[ICodeType]string
}

var (
	ctm  *errorCodeTypeMsgMap
	once sync.Once
)

// RegisterCode 注册常量代码
//
// Deprecated: RegisterCode is deprecated,it will be removed in the future.
//
// Please use RegisterCodeTable instead.
//
// i18nMsg key为错误码值，value为错误信息
//
// 与 RegisterCodeSingle 不同的是，此方法是将code表整个覆盖。
//
// 当code码重复时，后者覆盖前者
//
// # 此方法适用于注入【固定】的错误码表的场景
func RegisterCode(language LanguageCode, i18nMsg map[ICodeType]string) {
	ctm.mux.Lock()
	ctm.maps[language] = i18nMsg
	ctm.mux.Unlock()
}

// RegisterCodeTable 注册错误码表
//
// i18nMsg key为错误码值，value为错误信息
//
// 与 RegisterCodeSingle 不同的是，此方法是覆盖对应语言的整个code表。
//
// # 此方法适用于【固定】注入错误码表或需要覆盖默认错误码表的场景
func RegisterCodeTable(language LanguageCode, i18nMsg map[ICodeType]string) {
	ctm.mux.Lock()
	ctm.maps[language] = i18nMsg
	ctm.mux.Unlock()
}

// RegisterCodeSingle 注册单个错误码
//
// 与 RegisterCodeTable 不同的是，此方法更细粒度的注入错误码及对应的错误信息而不是整个覆盖。
//
// 当code码重复时，后者覆盖前者
//
// # 此方法适用于【动态】注入单个错误码的场景
func RegisterCodeSingle(language LanguageCode, code ICodeType, msg string) {
	ctm.mux.Lock()
	if _, ok := ctm.maps[language]; !ok {
		ctm.maps[language] = make(map[ICodeType]string)
	}
	ctm.maps[language][code] = msg
	ctm.mux.Unlock()
}

func init() {
	once.Do(func() {
		ctm = &errorCodeTypeMsgMap{
			mux:  &sync.RWMutex{},
			maps: make(map[LanguageCode]map[ICodeType]string),
		}

		for language, codeMsg := range initErrorCodeMsgMap {
			RegisterCodeTable(language, codeMsg)
		}
	})
}
