package api

// CookieStd cookie field definitions
type CookieStd struct {
	Name     string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}
