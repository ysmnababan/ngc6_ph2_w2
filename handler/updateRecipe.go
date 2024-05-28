package handler

import (
	"encoding/json"
	"net/http"
	"siang/model"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (repo *MysqlDB) UpdateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))

	var recipe model.Recipe

	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		http.Error(w, "error when getting body", http.StatusBadRequest)
		return
	}

	_, err = repo.DB.Exec("UPDATE recipes SET nama = ?, description =?, rating=? WHERE id= ?", recipe.Name, recipe.Description, recipe.Rating, id)
	if err != nil {
		http.Error(w, "error query", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("update successfully")
	w.WriteHeader(http.StatusOK)
}
