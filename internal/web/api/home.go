package api

import (
	"encoding/json"
	"net/http"

	"github.com/laetificat/slogger/pkg/slogger"
	"github.com/spf13/viper"

	"github.com/julienschmidt/httprouter"
)

/*
RegisterHomeHandler registers the home handler.
*/
func RegisterHomeHandler(router *httprouter.Router) {
	router.GET("/", Index)
}

/*
Index returns the API version.
*/
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := json.Marshal(struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}{"Pricewatcher API server", viper.GetString("version")})

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	header := w.Header()
	header.Set("Content-Type", "application/json")

	_, err = w.Write(body)
	if err != nil {
		slogger.Error(err.Error())
	}
}
