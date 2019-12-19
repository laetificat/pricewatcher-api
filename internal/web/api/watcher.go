package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/laetificat/pricewatcher/internal/helper"
	"github.com/laetificat/pricewatcher/internal/model"
	"github.com/laetificat/pricewatcher/internal/watcher"
	"github.com/laetificat/slogger/pkg/slogger"
)

/*
RegisterWatcherHandler registers the watcher handler.
*/
func RegisterWatcherHandler(router *httprouter.Router) {
	router.GET("/watchers", ListAll)
	router.GET("/watchers/run/:id", RunAll)
	router.GET("/watchers/delete/:id", DeleteOne)
	router.GET("/watchers/create", AddOne)
}

/*
ListAll returns a list of all the watchers, filters the list by url or domain if given as query params.
*/
func ListAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryValues := r.URL.Query()

	queryKeys := map[string]string{}

	if urlParam := queryValues.Get("url"); urlParam != "" {
		queryKeys["Url"] = urlParam
	}

	if domainParam := queryValues.Get("domain"); domainParam != "" {
		queryKeys["Domain"] = domainParam
	}

	var priceHistories []model.Watcher
	var err error

	priceHistories, err = watcher.List(queryKeys)

	if err != nil {
		slogger.Error(err.Error())
		return
	}

	jbody, err := json.Marshal(priceHistories)
	if err != nil {
		slogger.Error(err.Error())
		return
	}

	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(jbody)
	if err != nil {
		slogger.Error(err.Error())
	}
}

/*
RunAll registers all the jobs for the watchers in all the queues, if given an id it will only register watchers for the
queue with the given id.
*/
func RunAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ParamsID := p.ByName("id")

	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Access-Control-Allow-Origin", "*")

	if ParamsID == "" {
		err := watcher.RunAll()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			slogger.Error(err.Error())
			return
		}

		return
	}

	iID, err := strconv.Atoi(ParamsID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slogger.Error(err.Error())
		return
	}

	err = watcher.Run(iID)
	if err != nil {
		if strings.EqualFold(err.Error(), "key not found") {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			slogger.Info(err.Error())
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slogger.Error(err.Error())
		return
	}
}

/*
DeleteOne deleted a single watcher from the database based on given id.
*/
func DeleteOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Access-Control-Allow-Origin", "*")

	if id == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	iID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		slogger.Error(err.Error())
		return
	}

	err = watcher.Remove(iID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slogger.Error(err.Error())
	}
}

/*
AddOne registers a new watcher based on the given query parameters url and domain. It is possible to omit domain as
this will be added automatically.
*/
func AddOne(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryValues := r.URL.Query()
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Access-Control-Allow-Origin", "*")

	givenURL := queryValues.Get("url")
	givenDomain := queryValues.Get("domain")

	if givenDomain == "" {
		var err error
		givenDomain, err = helper.GuessDomain(givenURL)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			slogger.Error(err.Error())
			return
		}
	}

	if helper.IsSupported(givenDomain) {
		if err := watcher.Add(givenDomain, givenURL); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			slogger.Error(err.Error())
			return
		}
	}

	errorTxt := fmt.Sprintf("Given domain '%s' is not supported.", givenDomain)
	slogger.Info(errorTxt)
	http.Error(w, errorTxt, http.StatusNotAcceptable)
}
