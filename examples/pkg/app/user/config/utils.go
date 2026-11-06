package config

import "time"

//TestLoadTimeZonePanicWhileFailure tests loading the time zone
//
//Panic if loading fails
func TestLoadTimeZonePanicWhileFailure() {
	_, err := time.LoadLocation(GetGlobalConfig().Timezone)
	if err != nil {
		panic(err)
	}
}
