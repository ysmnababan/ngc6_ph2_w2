package handler

import (
	"encoding/json"
	"net/http"
	"siang/model"

	"github.com/julienschmidt/httprouter"
)

func (repo *MysqlDB) CreateRecipe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var recipe model.Recipe

	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		http.Error(w, "error when getting body", http.StatusBadRequest)
		return
	}

	result, err := repo.DB.Exec("INSERT INTO recipes (nama, description, rating) VALUES (?,?,?)", recipe.Name, recipe.Description, recipe.Rating)
	if err != nil {
		http.Error(w, "error query", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "error insert id ", http.StatusInternalServerError)
		return
	}
	recipe.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
	w.WriteHeader(http.StatusOK)
}
