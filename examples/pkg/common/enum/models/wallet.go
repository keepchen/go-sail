package models

type WalletStatusCode int

const (
	WalletStatusCodeNormal    = WalletStatusCode(0) //正常
	WalletStatusCodeForbidden = WalletStatusCode(1) //禁用
)

var walletStatusCodeMap = map[WalletStatusCode]string{
	WalletStatusCodeNormal:    "normal",
	WalletStatusCodeForbidden: "forbidden",
}

func (wsc WalletStatusCode) Int() int {
	return int(wsc)
}

func (wsc WalletStatusCode) String() string {
	return walletStatusCodeMap[wsc]
}
