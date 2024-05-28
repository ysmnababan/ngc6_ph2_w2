package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"siang/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(u model.User) (map[string]string, error) {
	response := map[string]string{}
	//create payload
	payload := jwt.MapClaims{
		"id":       u.Id,
		"email":    u.Email,
		"fullname": u.Fullname,
		"role":     u.Role,
	}

	// define the method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	//get secret code from .env
	err := godotenv.Load()
	if err != nil {
		return response, fmt.Errorf("error while getting secret key")
	}

	// get token string
	key := os.Getenv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return response, fmt.Errorf("error generating token: %v", err)
	}

	response["token"] = tokenString
	return response, nil
}

func (repo *MysqlDB) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var credential model.User

	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		http.Error(w, "error when getting body", http.StatusBadRequest)
		return
	}

	var storedUser model.User

	err = repo.DB.QueryRow("SELECT id, fullname, email, password, age, occupation, role FROM users WHERE email = ?", credential.Email).Scan(&storedUser.Id, &storedUser.Fullname, &storedUser.Email, &storedUser.Password, &storedUser.Age, &storedUser.Occupation, &storedUser.Role)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(credential.Password))
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// login succeed
	fmt.Println("login succeed")
	token, err := GenerateToken(storedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}
