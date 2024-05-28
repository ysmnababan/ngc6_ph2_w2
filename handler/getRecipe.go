package handler

import (
	"encoding/json"
	"net/http"
	"siang/model"

	"github.com/julienschmidt/httprouter"
)

func (repo *MysqlDB) GetRecipe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var recipes []model.Recipe

	rows, err := repo.DB.Query("SELECT id, nama, description, rating FROM recipes")
	if err != nil {
		http.Error(w, "error query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r model.Recipe
		err = rows.Scan(&r.ID, &r.Name, &r.Description, &r.Rating)
		if err != nil {
			http.Error(w, "error scan", http.StatusInternalServerError)
			return
		}

		recipes = append(recipes, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
	w.WriteHeader(http.StatusOK)

}
