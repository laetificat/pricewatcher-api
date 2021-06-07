package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/laetificat/pricewatcher-api/internal/log"
	"github.com/laetificat/pricewatcher-api/internal/model"
	"github.com/laetificat/pricewatcher-api/internal/queue"
)

/*
RegisterQueueHandler registers the queue handler.
*/
func RegisterQueueHandler(router *httprouter.Router) {
	router.GET("/queues/:name/next", GetNextItem)
	router.GET("/queues", GetAvailableQueues)
	router.GET("/queues/:name", GetQueueItems)
	router.POST("/queues/:name/add", AddQueueItem)
}

/*
AddQueueItem accepts a JSON encoded watcher model and adds it to the database.
*/
func AddQueueItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	queueName := p.ByName("name")

	if queueName == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	responseModel := model.Watcher{}
	err := json.NewDecoder(r.Body).Decode(&responseModel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	err = queue.Add(queueName, &responseModel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
}

/*
GetQueueItems returns a list of watcher models for a specific queue.
*/
func GetQueueItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	queueName := p.ByName("name")

	if queueName == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	watchers, err := queue.Get(queueName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	responseModel := struct {
		Jobs []*model.Watcher `json:"jobs"`
	}{watchers}

	responseBody, err := json.Marshal(responseModel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
}

/*
GetAvailableQueues returns a list of queue names that are registered.
*/
func GetAvailableQueues(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var queueNames []string
	for k := range queue.ListQueues() {
		queueNames = append(queueNames, k)
	}

	responseBody := struct {
		Queues []string `json:"queues"`
	}{queueNames}

	response, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
}

/*
GetNextItem returns the next item in the queue with the given name, returns nil if none is found.
*/
func GetNextItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	queueName := p.ByName("name")
	if queueName == "" {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}

	watcher, err := queue.Next(queueName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	responseBody, err := json.Marshal(watcher)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
}
