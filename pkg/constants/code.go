package constants

import "sync"

// ICodeType 错误码类型接口
type ICodeType interface {
	String() string
	Int() int
}

// CodeType 错误码类型
type CodeType int

type errorCodeTypeMsgMap struct {
	mux  *sync.RWMutex
	maps map[CodeType]string
}

var ctm *errorCodeTypeMsgMap

// RegisterCode 注册常量代码
//
// 当code码重复时，后者覆盖前者
func RegisterCode(code CodeType, msg string) {
	ctm.mux.Lock()
	ctm.maps[code] = msg
	ctm.mux.Unlock()
}

func init() {
	(&sync.Once{}).Do(func() {
		ctm = &errorCodeTypeMsgMap{
			mux:  &sync.RWMutex{},
			maps: make(map[CodeType]string),
		}

		for code, msg := range initErrorCodeMsgMap {
			RegisterCode(code, msg)
		}
	})
}
