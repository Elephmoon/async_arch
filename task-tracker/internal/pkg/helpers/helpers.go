package helpers

import (
	"log"
	"net/http"
)

func WriteHttpResponse(resp http.ResponseWriter, status int, payload []byte) {
	_, err := resp.Write(payload)
	if err != nil {
		log.Printf("cant write http response %v", err)
	}
	resp.WriteHeader(status)
}
