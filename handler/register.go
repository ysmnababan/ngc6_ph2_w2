package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"siang/model"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type MysqlDB struct {
	DB *sql.DB
}

func isValidEmail(email string) bool {
	// Define the regex pattern for a valid email address
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern
	re := regexp.MustCompile(emailRegexPattern)

	// Match the email with the regex pattern
	return re.MatchString(email)
}

func validateUser(u model.User) error {

	if u.Email == "" || u.Password == "" || u.Fullname == "" || u.Occupation == "" || u.Age <= 0 || u.Role == "" {
		return fmt.Errorf("invalid or missing parameter")
	}

	if !isValidEmail(u.Email) {
		return fmt.Errorf("email format is not valid")
	}

	if len(u.Password) < 8 {
		return fmt.Errorf("password length should be longer than 8 character")
	}

	if len(u.Fullname) < 6 || len(u.Fullname) > 15 {
		return fmt.Errorf("full name length must be 6 to 15 characters")
	}

	if u.Age < 17 {
		return fmt.Errorf("age must be 17 and above")
	}

	if !(u.Role == "admin" || u.Role == "superadmin") {
		return fmt.Errorf("role is invalid")
	}

	return nil
}

func (repo *MysqlDB) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var u model.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "error when getting body", http.StatusBadRequest)
		return
	}

	err = validateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedpwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing", http.StatusInternalServerError)
		return
	}

	result, err := repo.DB.Exec("INSERT INTO users (fullname, email, password, age, occupation,role) VALUES (?,?,?,?,?,?)", u.Fullname, u.Email, hashedpwd, u.Age, u.Occupation, u.Role)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "unable to get lastID", http.StatusInternalServerError)
		return
	}

	u.Id = int(lastID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}
