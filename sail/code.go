package sail

import "github.com/keepchen/go-sail/v3/constants"

// ICode 错误码接口
type ICode interface {
	// Register 将错误码信息注册到错误码表中
	//
	// - language，语言代码，需要符合 constants.LanguageCode 表中的值，例如en,en-US
	//
	// - code，错误码值
	//
	// - msg，错误提示信息
	Register(language string, code int, msg string)
}

type codeImpl struct{}

var ci ICode = &codeImpl{}

// Code 获取错误码接口方法
func Code() ICode {
	return ci
}

func (*codeImpl) Register(language string, code int, msg string) {
	constants.RegisterCodeSingle(constants.LanguageCode(language), constants.CodeType(code), msg)
}
