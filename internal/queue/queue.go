package queue

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/laetificat/pricewatcher/internal/model"
)

var queues = map[string]*list.List{}

/*
ListQueues returns a list of queues that are currently registered.
*/
func ListQueues() map[string]*list.List {
	return queues
}

/*
Get returns a list of items in a queue with the given name.
*/
func Get(queueName string) ([]*model.Watcher, error) {
	var watchers []*model.Watcher
	if queue, ok := queues[queueName]; ok {
		for e := queue.Front(); e != nil; e = e.Next() {
			watchers = append(watchers, e.Value.(*model.Watcher))
		}

		return watchers, nil
	}

	return nil, fmt.Errorf("could not find queue '%s'", queueName)
}

/*
Add adds the given watcher to the queue with the given name.
*/
func Add(queueName string, watcher *model.Watcher) error {
	if queue, ok := queues[queueName]; ok {
		for e := queue.Front(); e != nil; e = e.Next() {
			if e.Value.(*model.Watcher).ID == watcher.ID {
				return nil
			}
		}

		queue.PushBack(watcher)
		return nil
	}

	return fmt.Errorf("could not find queue '%s'", queueName)
}

/*
Next returns the first item from the queue the front with the given name, when returning it also removes it from the queue.
*/
func Next(name string) (*model.Watcher, error) {
	if queue, ok := queues[name]; ok {
		if queue.Front() != nil {
			queueElement := queue.Front()
			queue.Remove(queueElement)
			return queueElement.Value.(*model.Watcher), nil
		}

		return nil, nil
	}

	return nil, fmt.Errorf("could not find queue with name '%s'", name)
}

/*
GetNameForDomain returns a queue equivalent name for the given domain.
It prepends "queue_" to the given domain and replaces "." with "_".
*/
func GetNameForDomain(domain string) string {
	return "queue_" + strings.ReplaceAll(domain, ".", "_")
}

/*
Create creates a new queue with the given name if it does not exist.
*/
func Create(domain ...string) error {
	for _, v := range domain {
		if _, ok := queues[v]; ok {
			return fmt.Errorf("queue with name '%s' already exists", v)
		}
		queues[v] = list.New()
	}

	return nil
}
