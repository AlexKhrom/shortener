package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func LofInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("\n\n\n=====================")
		fmt.Println("time: ", time.Now().String()[:27])
		fmt.Println("r.Method: ", r.Method)
		//fmt.Println("vers method = ", r.Proto)
		fmt.Println("url: ", r.URL.Path)
		fmt.Println("=====================")

		next.ServeHTTP(w, r)
	})
}
