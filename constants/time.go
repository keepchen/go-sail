package constants

// @see https://timezonedb.com/time-zones
// Linux shell: timedatectl list-timezones
const (
	DefaultTimeZone   string = "Asia/Shanghai" //默认时区
	TimeZoneUTCSub11  string = "Pacific/Niue"
	TimeZoneUTCSub10  string = "Pacific/Honolulu"
	TimeZoneUTCSub9   string = "Pacific/Gambier"
	TimeZoneUTCSub8   string = "Asia/Brunei"
	TimeZoneUTCSub7   string = "America/Hermosillo"
	TimeZoneUTCSub6   string = "America/Belize"
	TimeZoneUTCSub5   string = "America/Eirunepe"
	TimeZoneUTCSub4   string = "America/Anguilla"
	TimeZoneUTCSub3   string = "Antarctica/Rothera"
	TimeZoneUTCSub2   string = "America/Noronha"
	TimeZoneUTCSub1   string = "Atlantic/Cape_Verde"
	TimeZoneUTC0      string = "Africa/Ouagadougou"
	TimeZoneUTCPlus1  string = "Africa/Douala"
	TimeZoneUTCPlus2  string = "Africa/Lubumbashi"
	TimeZoneUTCPlus3  string = "Antarctica/Syowa"
	TimeZoneUTCPlus4  string = "Asia/Yerevan"
	TimeZoneUTCPlus5  string = "Indian/Kerguelen"
	TimeZoneUTCPlus6  string = "Asia/Almaty"
	TimeZoneUTCPlus7  string = "Asia/Vientiane"
	TimeZoneUTCPlus8  string = "Asia/Shanghai"
	TimeZoneUTCPlus9  string = "Asia/Jayapura"
	TimeZoneUTCPlus10 string = "Pacific/Chuuk"
	TimeZoneUTCPlus11 string = "Pacific/Kosrae"
	TimeZoneUTCPlus12 string = "Pacific/Nauru"
)

const (
	DateLayout       string = "2006-01-02"           //日期格式
	TimeLayout       string = "15:04:05"             //时间格式
	DatetimeLayout   string = "2006-01-02 15:04:05"  //日期时间格式
	DatetimeTZLayout string = "2006-01-02T15:04:05Z" //日期时间格式-tz
)
