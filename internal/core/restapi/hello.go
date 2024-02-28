package restapi

import (
	"fmt"
	"net/http"
)

func sayHello(msg string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello, %s!\"}", msg)))
	}
}
