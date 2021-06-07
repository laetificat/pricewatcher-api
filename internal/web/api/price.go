package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/laetificat/pricewatcher-api/internal/model"
	"github.com/laetificat/pricewatcher-api/internal/watcher"
)

/*
RegisterPriceHandler registers the price handler.
*/
func RegisterPriceHandler(router *httprouter.Router) {
	router.POST("/prices/update/:id", UpdatePrice)
}

/*
UpdatePrice accepts a JSON encoded update model and uses that to update the watcher model's price in the database.
*/
func UpdatePrice(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	updateModel := model.Update{}
	if err := json.NewDecoder(r.Body).Decode(&updateModel); err != nil {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}

	err := watcher.Update(&updateModel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
