package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
)

func TestSetupOption(t *testing.T) {
	t.Run("SetupOption-Panic", func(t *testing.T) {
		opt := Option{
			Timezone: "Earth/Unknown",
		}
		assert.Panics(t, func() {
			SetupOption(opt)
		})
	})

	t.Run("SetupOption", func(t *testing.T) {
		emptyDataTypes := []int{DefaultEmptyDataStructNull, DefaultEmptyDataStructObject, DefaultEmptyDataStructArray, DefaultEmptyDataStructString, 999}
		values := []interface{}{nil, struct{}{}, []bool{}, "", nil}
		for idx, dt := range emptyDataTypes {
			opt := Option{
				Timezone:             constants.DefaultTimeZone,
				EmptyDataStruct:      dt,
				ForceHttpCode200:     true,
				DetectAcceptLanguage: true,
				FuncBeforeWrite: func(request *http.Request, entryAtUnixNano int64, requestId, spanId string, httpCode int, writeData dto.Base) {
					//do something...
				},

				ErrNoneCode:                      constants.CodeType(90000000),
				ErrRequestParamsInvalidCode:      constants.CodeType(90000001),
				ErrAuthorizationTokenInvalidCode: constants.CodeType(90000002),
				ErrInternalServerErrorCode:       constants.CodeType(90000003),
			}
			SetupOption(opt)

			assert.Equal(t, true, forceHttpCode200)
			assert.Equal(t, true, detectAcceptLanguage)
			assert.Equal(t, constants.DefaultTimeZone, timezone)
			assert.NotNil(t, loc)
			assert.NotNil(t, funcBeforeWrite)

			t.Log("idx", idx)
			assert.Equal(t, emptyDataField, values[idx])
		}
	})
}

func TestDefaultSetupOption(t *testing.T) {
	t.Run("DefaultSetupOption", func(t *testing.T) {
		t.Log(DefaultSetupOption())
	})
}
