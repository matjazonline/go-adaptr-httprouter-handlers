package adaptr_test

import (
	adaptr "github.com/matjazonline/go-httprouter-adapter-handlers"
	"testing"
)
import "net/http"
import "net/http/httptest"
import "fmt"

func TestCallOnce(t *testing.T) {
	println("START TEST 1")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("handler called")
		fmt.Fprint(w, "handler called")
	})

	initAdaptr := adaptr.CallOnce(func(w http.ResponseWriter, r *http.Request) {
		println("INIT CALLED ONCE")
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	baseAdaptrs := []adaptr.Adapter{initAdaptr, testAdaptr("-3"), testAdaptr("-2")}
	parentAdaptrs := append(baseAdaptrs, []adaptr.Adapter{initAdaptr, testAdaptr("-1"), testAdaptr("0")}...)

	adaptr.Adapt(handler, append(parentAdaptrs, initAdaptr, testAdaptr("1"), testAdaptr("2"))...).ServeHTTP(w, r)
	adaptr.Adapt(handler, initAdaptr, testAdaptr("3"), testAdaptr("4")).ServeHTTP(w, r)

}

func testAdaptr(logVal string) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			println("TestAdaptr called" + logVal)
			h.ServeHTTP(w, r)
		})
	}
}
