package sail

import "github.com/keepchen/go-sail/v3/constants"

// ICode error code interface
type ICode interface {
	// Register registers an error code message in the error code table
	//
	// - language, language code matching a value in constants.LanguageCode, e.g. en,en-US
	//
	// - code, error code value
	//
	// - msg, error message
	Register(language string, code int, msg string)
}

type codeImpl struct{}

var ci ICode = &codeImpl{}

// Code retrieves the error code interface
func Code() ICode {
	return ci
}

func (*codeImpl) Register(language string, code int, msg string) {
	constants.RegisterCodeSingle(constants.LanguageCode(language), constants.CodeType(code), msg)
}
