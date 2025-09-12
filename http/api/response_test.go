package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
)

// 创建gin测试下上文和引擎
func createTestContextAndEngine() (*gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()

	//创建测试用的Request（可自定义请求参数）
	req, _ := http.NewRequest("GET", "/test?name=foo", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", uuid.New().String())
	req.Header.Set("requestId", uuid.New().String())
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7,la;q=0.6")
	c, r := gin.CreateTestContext(w)
	c.Request = req
	c.Set("requestId", uuid.New().String())

	return c, r
}

func TestResponse(t *testing.T) {
	t.Run("Response", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		t.Log(Response(c))
	})
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		t.Log(New(c))
	})
}

type testerResponseData struct {
	dto.Base
	Data string `json:"data"`
}

func (v testerResponseData) GetData() interface{} {
	return v.Data
}

func TestBuilder(t *testing.T) {
	t.Run("Builder", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		t.Log(Response(c).Builder(constants.ErrNone, testerResponseData{}))
	})
}

func TestSuccess(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		Response(c).Success()
	})
}

func TestFailure(t *testing.T) {
	t.Run("Failure", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		Response(c).Failure()
		c2, _ := createTestContextAndEngine()
		Response(c2).Failure("oops~")
	})
}

func TestFailure200(t *testing.T) {
	t.Run("Failure200", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		Response(c).Failure200(constants.ErrRequestParamsInvalid)
	})
}

func TestFailure400(t *testing.T) {
	t.Run("Failure400", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		Response(c).Failure400(constants.ErrRequestParamsInvalid)
	})
}

func TestFailure500(t *testing.T) {
	t.Run("Failure500", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		Response(c).Failure500(constants.ErrRequestParamsInvalid)
	})
}

func TestData(t *testing.T) {
	t.Run("Data", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		Response(c).Data(nil)
		c2, _ := createTestContextAndEngine()
		anotherErrNoneCode = constants.CodeType(200)
		Response(c2).Data(nil)
	})
}

func TestBundle(t *testing.T) {
	t.Run("Bundle()", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		t.Log(Response(c).Bundle(200, testerResponseData{}))
	})
}

func TestWrap(t *testing.T) {
	t.Run("Wrap", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		t.Log(Response(c).Wrap(constants.ErrNone, testerResponseData{}))
	})
}

func TestSimpleAssemble(t *testing.T) {
	t.Run("SimpleAssemble", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		t.Log(Response(c).SimpleAssemble(constants.ErrNone, testerResponseData{}))
	})
}

func TestAssemble(t *testing.T) {
	t.Run("Assemble", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		t.Log(Response(c).Assemble(constants.ErrNone, testerResponseData{}))
	})
}

func TestStatus(t *testing.T) {
	t.Run("Status", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		t.Log(Response(c).Status(http.StatusNotFound))
	})
}

func TestSendWithCode(t *testing.T) {
	t.Run("SendWithCode", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		Response(c).Status(http.StatusNotFound).SendWithCode(http.StatusOK)
	})
}

func TestSend(t *testing.T) {
	t.Run("Send", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		Response(c).Status(http.StatusNotFound).Send()
	})
}

func TestMergeBody(t *testing.T) {
	t.Run("MergeBody-DetectLanguage", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		re := responseEngine{
			engine:    c,
			httpCode:  http.StatusOK,
			data:      nil,
			requestId: uuid.New().String(),
		}

		detectAcceptLanguage = true

		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		for _, code := range codes {
			t.Log(re.mergeBody(code, nil))
		}
	})

	t.Run("MergeBody-Code", func(t *testing.T) {
		c, _ := createTestContextAndEngine()
		re := responseEngine{
			engine:    c,
			httpCode:  http.StatusOK,
			data:      nil,
			requestId: uuid.New().String(),
		}

		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		for _, code := range codes {
			t.Log(re.mergeBody(code, nil))
		}
	})

	t.Run("MergeBody-XRequestId", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-Id", uuid.New().String())
		//req.Header.Set("requestId", uuid.New().String())
		//req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7,la;q=0.6")

		c.Request = req

		re := responseEngine{
			engine:    c,
			httpCode:  http.StatusOK,
			data:      nil,
			requestId: uuid.New().String(),
		}

		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		for _, code := range codes {
			t.Log(re.mergeBody(code, nil))
		}
	})

	t.Run("MergeBody-RequestId", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		//req.Header.Set("X-Request-Id", uuid.New().String())
		req.Header.Set("requestId", uuid.New().String())
		//req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7,la;q=0.6")

		c.Request = req

		re := responseEngine{
			engine:    c,
			httpCode:  http.StatusOK,
			data:      nil,
			requestId: uuid.New().String(),
		}

		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		for _, code := range codes {
			t.Log(re.mergeBody(code, nil))
		}
	})

	t.Run("MergeBody-RequestIdFromContext", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		c.Request = req
		c.Set("requestId", uuid.New().String())

		re := responseEngine{
			engine:    c,
			httpCode:  http.StatusOK,
			data:      nil,
			requestId: uuid.New().String(),
		}

		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		for _, code := range codes {
			t.Log(re.mergeBody(code, nil))
		}
	})

	t.Run("MergeBody-LocalIsNil", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		c.Request = req
		c.Set("requestId", uuid.New().String())

		loc = nil

		re := responseEngine{
			engine:    c,
			httpCode:  http.StatusOK,
			data:      nil,
			requestId: uuid.New().String(),
		}

		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		for _, code := range codes {
			t.Log(re.mergeBody(code, nil))
		}
	})

	t.Run("MergeBody-LocalNotNil", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		c.Request = req
		c.Set("requestId", uuid.New().String())

		timezone = constants.DefaultTimeZone

		re := responseEngine{
			engine:    c,
			httpCode:  http.StatusOK,
			data:      nil,
			requestId: uuid.New().String(),
		}

		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		for _, code := range codes {
			t.Log(re.mergeBody(code, nil))
		}
	})

	t.Run("MergeBody-AnotherErrorCode", func(t *testing.T) {
		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}
		statusCodes := []int{
			http.StatusOK,
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusInternalServerError,
		}

		forceHttpCode200 = false
		anotherErrNoneCode = constants.ErrNone
		anotherErrRequestParamsInvalidCode = constants.ErrRequestParamsInvalid
		anotherErrAuthorizationTokenInvalidCode = constants.ErrAuthorizationTokenInvalid
		anotherErrInternalServerErrorCode = constants.ErrInternalServerError

		for idx, code := range codes {
			c, _ := createTestContextAndEngine()

			re := responseEngine{
				engine:    c,
				data:      nil,
				requestId: uuid.New().String(),
			}

			t.Log(re.mergeBody(code, nil))
			t.Log(forceHttpCode200, anotherErrNoneCode, anotherErrInternalServerErrorCode, code, re.httpCode)
			assert.Equal(t, true, (code == 0 && re.httpCode == http.StatusOK) || (code != 0 && re.httpCode != http.StatusOK))
			assert.Equal(t, statusCodes[idx], re.httpCode)
		}
	})

	t.Run("MergeBody-ForceHttp200", func(t *testing.T) {
		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		forceHttpCode200 = true

		for _, code := range codes {
			c, _ := createTestContextAndEngine()

			re := responseEngine{
				engine:    c,
				data:      nil,
				requestId: uuid.New().String(),
			}

			t.Log(re.mergeBody(code, nil))
			assert.Equal(t, true, re.httpCode == http.StatusOK)
		}
	})

	t.Run("MergeBody-Messages", func(t *testing.T) {
		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		for _, code := range codes {
			c, _ := createTestContextAndEngine()

			re := responseEngine{
				engine:    c,
				httpCode:  http.StatusOK,
				data:      nil,
				requestId: uuid.New().String(),
			}

			t.Log(re.mergeBody(code, nil, "error1", "error2", "error3"))
		}
	})

	t.Run("MergeBody-DataInterface", func(t *testing.T) {
		codes := []constants.CodeType{
			constants.ErrNone,
			constants.ErrRequestParamsInvalid,
			constants.ErrAuthorizationTokenInvalid,
			constants.ErrInternalServerError,
		}

		forceHttpCode200 = true

		for _, code := range codes {
			c, _ := createTestContextAndEngine()

			re := responseEngine{
				engine:    c,
				httpCode:  http.StatusOK,
				data:      nil,
				requestId: uuid.New().String(),
			}
			t.Log(re.mergeBody(code, testerResponseData{Data: "123"}, "error1", "error2", "error3"))
			t.Log(re.mergeBody(code, nil, "error1", "error2", "error3"))
			t.Log(re.mergeBody(code, (*testerResponseData)(nil), "error1", "error2", "error3"))
			t.Log(re.mergeBody(code, &testerResponseData{Data: "abc"}, "error1", "error2", "error3"))
			t.Log(re.mergeBody(code, []string{"1", "2", "3"}, "error1", "error2", "error3"))
			t.Log(re.mergeBody(code, []string{}, "error1", "error2", "error3"))
		}
	})

	t.Run("MergeBody-EmptyDataField", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		re := responseEngine{
			engine:    c,
			httpCode:  http.StatusOK,
			data:      nil,
			requestId: uuid.New().String(),
		}

		forceHttpCode200 = true
		emptyDataTypes := []int{DefaultEmptyDataStructNull, DefaultEmptyDataStructObject, DefaultEmptyDataStructArray, DefaultEmptyDataStructString, 999}

		for _, dt := range emptyDataTypes {

			emptyDataField = dt

			t.Log(re.mergeBody(constants.ErrNone, testerResponseData{}, "error1", "error2", "error3"))
			t.Log(re.mergeBody(constants.ErrNone, nil, "error1", "error2", "error3"))
			t.Log(re.mergeBody(constants.ErrNone, testerResponseData{Data: "123"}, "error1", "error2", "error3"))
			t.Log(re.mergeBody(constants.ErrNone, (*testerResponseData)(nil), "error1", "error2", "error3"))
			t.Log(re.mergeBody(constants.ErrNone, &testerResponseData{Data: "abc"}, "error1", "error2", "error3"))
			t.Log(re.mergeBody(constants.ErrNone, []string{"1", "2", "3"}, "error1", "error2", "error3"))
			t.Log(re.mergeBody(constants.ErrNone, []string{}, "error1", "error2", "error3"))
		}
	})

	t.Run("MergeBody-FuncBeforeWrite", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		c.Set("requestId", uuid.New().String())
		c.Set("spanId", uuid.New().String())
		c.Set("entryAt", time.Now().UnixNano())

		emptyDataTypes := []int{DefaultEmptyDataStructNull, DefaultEmptyDataStructObject, DefaultEmptyDataStructArray, DefaultEmptyDataStructString, 999}
		fn := func(request *http.Request, entryAtUnixNano int64, requestId, spanId string, httpCode int, writeData dto.Base) {
			//do something...
		}

		for _, dt := range emptyDataTypes {

			SetupOption(Option{
				ForceHttpCode200: true,
				EmptyDataStruct:  dt,
				FuncBeforeWrite:  fn,
			})

			re := responseEngine{
				engine:    c,
				data:      nil,
				requestId: uuid.New().String(),
			}
			re.mergeBody(constants.ErrNone, testerResponseData{}, "error1", "error2", "error3").Send()
			re = responseEngine{
				engine:    c,
				data:      nil,
				requestId: uuid.New().String(),
			}
			re.mergeBody(constants.ErrNone, dto.Base{}, "error1", "error2", "error3").Send()
		}
	})
}

func TestIsTypedNil(t *testing.T) {
	t.Run("IsTypedNil", func(t *testing.T) {
		var mp map[string]string
		assert.Equal(t, true, isTypedNil(mp))
		mp = map[string]string{"name": "go-sail"}
		assert.Equal(t, false, isTypedNil(mp))

		var ch chan struct{}
		assert.Equal(t, true, isTypedNil(ch))
		ch = make(chan struct{}, 1)
		assert.Equal(t, false, isTypedNil(ch))

		var fn func()
		assert.Equal(t, true, isTypedNil(fn))
		fn = func() {
			fmt.Println("hello")
		}
		assert.Equal(t, false, isTypedNil(fn))

		var face interface{}
		assert.Equal(t, true, isTypedNil(face))
		face = 123
		assert.Equal(t, false, isTypedNil(face))

		var base *dto.Base
		assert.Equal(t, true, isTypedNil(base))
		base = &dto.Base{}
		assert.Equal(t, false, isTypedNil(base))

		var slice []string
		assert.Equal(t, true, isTypedNil(slice))
		slice = []string{"a"}
		assert.Equal(t, false, isTypedNil(slice))

		var ptr unsafe.Pointer
		assert.Equal(t, true, isTypedNil(ptr))
		ptr = unsafe.Pointer(&dto.Base{})
		assert.Equal(t, false, isTypedNil(ptr))
	})
}
