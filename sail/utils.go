package sail

import (
	"sync"

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

var _ IUtils = &utilsImpl{}

var utilsImplPool = sync.Pool{
	New: func() interface{} {
		return &utilsImpl{}
	},
}

// Utils 获取工具类方法
func Utils() IUtils {
	u := utilsImplPool.Get().(*utilsImpl)
	return u
}

func (u *utilsImpl) Aes() utils.IAes {
	utilsImplPool.Put(u)

	return utils.Aes()
}

func (u *utilsImpl) Base64() utils.IBase64 {
	utilsImplPool.Put(u)

	return utils.Base64()
}

func (u *utilsImpl) Cert() utils.ICert {
	utilsImplPool.Put(u)

	return utils.Cert()
}

func (u *utilsImpl) Crc32() utils.ICrc32 {
	utilsImplPool.Put(u)

	return utils.Crc32()
}

func (u *utilsImpl) Crc64() utils.ICrc64 {
	utilsImplPool.Put(u)

	return utils.Crc64()
}

func (u *utilsImpl) Datetime() utils.IDatetime {
	utilsImplPool.Put(u)

	return utils.Datetime()
}

func (u *utilsImpl) Domain() utils.IDomain {
	utilsImplPool.Put(u)

	return utils.Domain()
}

func (u *utilsImpl) File() utils.IFile {
	utilsImplPool.Put(u)

	return utils.File()
}

func (u *utilsImpl) Gzip() utils.IGzip {
	utilsImplPool.Put(u)

	return utils.Gzip()
}

func (u *utilsImpl) HttpClient() utils.IHttpClient {
	utilsImplPool.Put(u)

	return utils.HttpClient()
}

func (u *utilsImpl) IP() utils.IIP {
	utilsImplPool.Put(u)

	return utils.IP()
}

func (u *utilsImpl) MD5() utils.IMD5 {
	utilsImplPool.Put(u)

	return utils.MD5()
}

func (u *utilsImpl) Number() utils.INumber {
	utilsImplPool.Put(u)

	return utils.Number()
}

func (u *utilsImpl) RedisLocker(client ...redisLib.UniversalClient) utils.IRedisLocker {
	utilsImplPool.Put(u)

	return utils.RedisLocker(client...)
}

func (u *utilsImpl) RSA() utils.IRsa {
	utilsImplPool.Put(u)

	return utils.RSA()
}

func (u *utilsImpl) Signal() utils.ISignal {
	utilsImplPool.Put(u)

	return utils.Signal()
}

func (u *utilsImpl) SM4() utils.ISM4 {
	utilsImplPool.Put(u)

	return utils.SM4()
}

func (u *utilsImpl) String() utils.IString {
	utilsImplPool.Put(u)

	return utils.String()
}

func (u *utilsImpl) Swagger() utils.ISwagger {
	utilsImplPool.Put(u)

	return utils.Swagger()
}

func (u *utilsImpl) Validator() utils.IValidator {
	utilsImplPool.Put(u)

	return utils.Validator()
}

func (u *utilsImpl) Version() utils.IVersion {
	utilsImplPool.Put(u)

	return utils.Version()
}

func (u *utilsImpl) WebPush() utils.IWebPush {
	utilsImplPool.Put(u)

	return utils.WebPush()
}
