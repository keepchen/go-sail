package api

// CookieStd cookie字段定义
type CookieStd struct {
	Name     string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}
