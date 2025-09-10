package sail

import (
	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/utils"
)

// IUtils 工具类接口
type IUtils interface {
	Aes() utils.IAes
	Base64() utils.IBase64
	Cert() utils.ICert
	Crc32() utils.ICrc32
	Crc64() utils.ICrc64
	Datetime() utils.IDatetime
	Domain() utils.IDomain
	File() utils.IFile
	Gzip() utils.IGzip
	HttpClient() utils.IHttpClient
	IP() utils.IIP
	MD5() utils.IMD5
	Number() utils.INumber
	RedisLocker(client ...redisLib.UniversalClient) utils.IRedisLocker
	RSA() utils.IRsa
	Signal() utils.ISignal
	SM4() utils.ISM4
	String() utils.IString
	Swagger() utils.ISwagger
	Validator() utils.IValidator
	Version() utils.IVersion
	WebPush() utils.IWebPush
}

type utilsImpl struct{}

var ui IUtils = &utilsImpl{}

// Utils 获取工具类方法
func Utils() IUtils {
	return ui
}

func (u *utilsImpl) Aes() utils.IAes {
	return utils.Aes()
}

func (u *utilsImpl) Base64() utils.IBase64 {
	return utils.Base64()
}

func (u *utilsImpl) Cert() utils.ICert {
	return utils.Cert()
}

func (u *utilsImpl) Crc32() utils.ICrc32 {
	return utils.Crc32()
}

func (u *utilsImpl) Crc64() utils.ICrc64 {
	return utils.Crc64()
}

func (u *utilsImpl) Datetime() utils.IDatetime {
	return utils.Datetime()
}

func (u *utilsImpl) Domain() utils.IDomain {
	return utils.Domain()
}

func (u *utilsImpl) File() utils.IFile {
	return utils.File()
}

func (u *utilsImpl) Gzip() utils.IGzip {
	return utils.Gzip()
}

func (u *utilsImpl) HttpClient() utils.IHttpClient {
	return utils.HttpClient()
}

func (u *utilsImpl) IP() utils.IIP {
	return utils.IP()
}

func (u *utilsImpl) MD5() utils.IMD5 {
	return utils.MD5()
}

func (u *utilsImpl) Number() utils.INumber {
	return utils.Number()
}

func (u *utilsImpl) RedisLocker(client ...redisLib.UniversalClient) utils.IRedisLocker {
	return utils.RedisLocker(client...)
}

func (u *utilsImpl) RSA() utils.IRsa {
	return utils.RSA()
}

func (u *utilsImpl) Signal() utils.ISignal {
	return utils.Signal()
}

func (u *utilsImpl) SM4() utils.ISM4 {
	return utils.SM4()
}

func (u *utilsImpl) String() utils.IString {
	return utils.String()
}

func (u *utilsImpl) Swagger() utils.ISwagger {
	return utils.Swagger()
}

func (u *utilsImpl) Validator() utils.IValidator {
	return utils.Validator()
}

func (u *utilsImpl) Version() utils.IVersion {
	return utils.Version()
}

func (u *utilsImpl) WebPush() utils.IWebPush {
	return utils.WebPush()
}
