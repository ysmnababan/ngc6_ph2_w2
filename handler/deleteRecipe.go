package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (repo *MysqlDB) DeleteRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))

	_, err := repo.DB.Exec("DELETE FROM recipes WHERE id = ?", id)
	if err != nil {
		http.Error(w, "error query", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("delete successfully")
	w.WriteHeader(http.StatusOK)
}
