package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"siang/config"
	"siang/handler"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func logRequest(message string) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			log.Println(message)
			next(w, r, p)
		}
	}
}
func AuthMiddleware(next httprouter.Handle, requireSuperUser bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// get token from header
		tokenString := r.Header.Get("Auth")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//get secret code from .env
		err := godotenv.Load()
		if err != nil {
			fmt.Println("error while getting secret key")
			return
		}

		// get token string
		key := os.Getenv("SECRET_KEY")

		// get token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(key), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// change token -> struct
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Error parsing claims")
			http.Error(w, "error getting claim", http.StatusInternalServerError)
			return
		}

		fmt.Println(claims["role"])
		if requireSuperUser {
			if claims["role"] != "superadmin" {
				http.Error(w, "only for superuser", http.StatusUnauthorized)
				return
			}
		}

		next(w, r, p)
	}
}

func main() {
	db := config.Connect()
	defer db.Close()

	handler := &handler.MysqlDB{DB: db}
	router, server := config.InitServer()
	router.POST("/register", logRequest("request sent to POST /register")(handler.Register))
	router.POST("/login", logRequest("request sent to POST /login")(handler.Login))

	// auth needed
	router.GET("/recipes", AuthMiddleware(handler.GetRecipe, false))
	router.POST("/recipe", AuthMiddleware(handler.CreateRecipe, true))
	router.PUT("/recipe/:id", AuthMiddleware(handler.UpdateRecipe, false))
	router.DELETE("/recipe/:id", AuthMiddleware(handler.DeleteRecipe, true))

	fmt.Println("server running on localhost:8080")
	panic(server.ListenAndServe())
}
