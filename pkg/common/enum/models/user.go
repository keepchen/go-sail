package models

type UserStatusCode int

const (
	UserStatusCodeNormal    = UserStatusCode(0) //正常
	UserStatusCodeForbidden = UserStatusCode(1) //禁用
)

var userStatusCodeMap = map[UserStatusCode]string{
	UserStatusCodeNormal:    "normal",
	UserStatusCodeForbidden: "forbidden",
}

func (usc UserStatusCode) Int() int {
	return int(usc)
}

func (usc UserStatusCode) String() string {
	return userStatusCodeMap[usc]
}
