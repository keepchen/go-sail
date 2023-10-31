package config

import "time"

//TestLoadTimeZonePanicWhileFailure 测试加载时区
//
//如果加载失败，将panic
func TestLoadTimeZonePanicWhileFailure() {
	_, err := time.LoadLocation(GetGlobalConfig().Timezone)
	if err != nil {
		panic(err)
	}
}
