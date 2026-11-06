package vo

import (
	"github.com/keepchen/go-sail/v3/constants"
)

// IRequest common request interface
type IRequest interface {
	//Validator validates request parameters
	Validator() (constants.ICodeType, error)
}
