package models

type UserStatusCode int

const (
	UserStatusCodeNormal    = UserStatusCode(0) //Normal
	UserStatusCodeForbidden = UserStatusCode(1) //Disabled
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
