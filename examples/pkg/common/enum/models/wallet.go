package models

type WalletStatusCode int

const (
	WalletStatusCodeNormal    = WalletStatusCode(0) //Normal
	WalletStatusCodeForbidden = WalletStatusCode(1) //Disabled
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
