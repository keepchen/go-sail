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
// i18nMsg key为错误码值，value为错误信息
//
// 当code码重复时，后者覆盖前者
func RegisterCode(language LanguageCode, i18nMsg map[ICodeType]string) {
	ctm.mux.Lock()
	ctm.maps[language] = i18nMsg
	ctm.mux.Unlock()
}

func init() {
	once.Do(func() {
		ctm = &errorCodeTypeMsgMap{
			mux:  &sync.RWMutex{},
			maps: make(map[LanguageCode]map[ICodeType]string),
		}

		for language, msg := range initErrorCodeMsgMap {
			RegisterCode(language, msg)
		}
	})
}
