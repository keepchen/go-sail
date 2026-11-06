package constants

import "sync"

// ICodeType error code interface
type ICodeType interface {
	String(language ...string) string
	Int() int
}

// CodeType error code type
type CodeType int

// MMBox language code-to-error code mapping type
type MMBox map[LanguageCode]map[ICodeType]string

// errorCodeTypeMsgMap uses sync.Map for concurrency safety
// Outer map: LanguageCode -> *sync.Map
// Inner map: ICodeType -> string
type errorCodeTypeMsgMap struct {
	maps sync.Map // map[LanguageCode]*sync.Map
}

var (
	ctm  *errorCodeTypeMsgMap
	once sync.Once
)

// RegisterCodeTable registers an error code table
//
// i18nMsg keys are error codes and values are error messages
//
// Unlike RegisterCodeSingle, this method replaces the entire code table for the given language.
//
// # Use this method for [static] error code table registration or to override the default table
func RegisterCodeTable(language LanguageCode, i18nMsg map[ICodeType]string) {
	m := &sync.Map{}
	for k, v := range i18nMsg {
		m.Store(k, v)
	}
	ctm.maps.Store(language, m)
}

// RegisterCodeSingle registers a single error code
//
// Unlike RegisterCodeTable, this method registers individual error codes and messages without replacing the entire table.
//
// For duplicate codes, the later value overwrites the earlier one
//
// # Use this method for [dynamic] registration of individual error codes
func RegisterCodeSingle(language LanguageCode, code ICodeType, msg string) {
	// Ensure the inner map exists
	val, _ := ctm.maps.LoadOrStore(language, &sync.Map{})
	inner := val.(*sync.Map)
	inner.Store(code, msg)
}

// GetCodeMsg returns the message for an error code
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
