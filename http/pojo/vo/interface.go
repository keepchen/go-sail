package vo

import (
	"github.com/keepchen/go-sail/v3/constants"
)

// IRequest 统一请求接口
type IRequest interface {
	//Validator 校验请求参数
	Validator() (constants.ICodeType, error)
}
