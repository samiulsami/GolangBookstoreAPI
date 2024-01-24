package Auth

import "net/http"

func GetJWTToken(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("asdasd"))
}
