package middlewares

import (
	"errors"
	"github.com/DaggerJackfast/gst/src/controllers"
	"github.com/DaggerJackfast/gst/src/token"
	"net/http"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		err := token.TokenValid(r)
		if err != nil {
			controllers.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized. Token is not valid."))
			return
		}
		next(w, r)
	}
}
