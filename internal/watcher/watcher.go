package watcher

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/laetificat/pricewatcher-api/internal/log"
	"github.com/laetificat/pricewatcher-api/internal/model"
	"github.com/laetificat/pricewatcher-api/internal/queue"
	"github.com/spf13/viper"

	bolt "go.etcd.io/bbolt"
)

// SupportedDomains is the list of supported domains.
var SupportedDomains = []string{
	"bol.com",
	"ebay.nl",
	"coolblue.nl",
}

/*
Add registers a new watcher object in the database with the given domain and url.
*/
func Add(domain, url string) error {
	db, err := bolt.Open(viper.GetString("database_file"), 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	watcher := model.Watcher{
		URL:          url,
		Domain:       domain,
		IsChecking:   false,
		PriceHistory: []model.Price{},
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, inErr := tx.CreateBucketIfNotExists([]byte("watchers"))
		if inErr != nil {
			return inErr
		}

		id, _ := b.NextSequence()
		watcher.ID = int(id)

		w, inErr := json.Marshal(watcher)
		if inErr != nil {
			return inErr
		}

		return b.Put(itob(watcher.ID), w)
	})

	return err
}

/*
List returns all watcher models from the database, filters items based on a map.

Example:
List(map[string]string{"domain": "bol.com"})
*/
func List(keys map[string]string) ([]model.Watcher, error) {
	watcherList := []model.Watcher{}

	db, err := bolt.Open(viper.GetString("database_file"), 0600, nil)
	if err != nil {
		return watcherList, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, inErr := tx.CreateBucketIfNotExists([]byte("watchers"))
		if inErr != nil {
			return inErr
		}

		inErr = b.ForEach(func(bk, bv []byte) error {
			watcher := model.Watcher{}
			inInErr := json.Unmarshal(bv, &watcher)
			if inInErr != nil {
				return inInErr
			}

			watcherStruct := structs.Map(watcher)

			if len(keys) == 0 {
				watcherList = append(watcherList, watcher)
			} else {
				for k, v := range keys {
					if val, ok := watcherStruct[k]; ok {
						if val == v {
							watcherList = append(watcherList, watcher)
						}
					}
				}
			}

			return nil
		})

		return inErr
	})

	return watcherList, err
}

/*
Remove removes a watcher model from the database based on ID.
*/
func Remove(id int) error {
	db, err := bolt.Open(viper.GetString("database_file"), 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("watchers"))
		if err != nil {
			return err
		}
		c := b.Cursor()

		for k, _ := c.Seek(itob(id)); k != nil && !bytes.Equal(k, itob(id)); k, _ = c.Next() {
			if err := c.Delete(); err != nil {
				return err
			}
		}

		return nil
	})
}

/*
RemoveAll removes all the registered watchers.
*/
func RemoveAll() error {
	db, err := bolt.Open(viper.GetString("database_file"), 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("watchers"))
		if err != nil {
			return err
		}

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			err := b.Delete(k)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

/*
Run adds a single watcher from the database to the queue as a job based on ID.
*/
func Run(id int) error {
	db, err := bolt.Open(viper.GetString("database_file"), 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("watchers"))
		if err != nil {
			return err
		}

		c := b.Cursor()
		k, v := c.Seek(itob(id))

		if k == nil || !bytes.Equal(k, itob(id)) {
			return fmt.Errorf("key not found")
		}

		client := &http.Client{}
		return addToQueue(client, v)
	})
}

/*
RunAll adds all the watchers to the queue as a job from the database.
*/
func RunAll() error {
	db, err := bolt.Open(viper.GetString("database_file"), 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("watchers"))
		if err != nil {
			return err
		}

		client := &http.Client{}

		err = b.ForEach(func(k, v []byte) error {
			return addToQueue(client, v)
		})

		return err
	})
}

/*
Update adds the given price from the update model to the watcher model that is found with the update model id.
*/
func Update(updateModel *model.Update) error {
	db, err := bolt.Open(viper.GetString("database_file"), 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("watchers"))
		if err != nil {
			return err
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if !bytes.Equal(k, itob(updateModel.ID)) {
				continue
			}

			bWatcher := model.Watcher{}
			if err := json.Unmarshal(v, &bWatcher); err != nil {
				return err
			}

			bWatcher.Name = updateModel.Name
			bWatcher.LastChecked = updateModel.Price.Timestamp
			bWatcher.PriceHistory = append(bWatcher.PriceHistory, updateModel.Price)

			w, err := json.Marshal(bWatcher)
			if err != nil {
				return err
			}

			return b.Put(k, w)
		}

		return nil
	})
}

func addToQueue(client *http.Client, v []byte) error {
	watcher := model.Watcher{}
	err := json.Unmarshal(v, &watcher)
	if err != nil {
		return err
	}

	if time.Since(watcher.LastChecked).Hours() <= viper.GetFloat64("watcher.check_interval") {
		return nil
	}

	log.Debug(fmt.Sprintf("Adding item to queue '%s'", watcher.Domain))
	res, err := client.Post(
		viper.GetString("webserver.address")+"/queues/"+queue.GetNameForDomain(watcher.Domain)+"/add",
		"application/json",
		bytes.NewBuffer(v),
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode > 200 {
		return fmt.Errorf("adding job failed, stats code %s", strconv.Itoa(res.StatusCode))
	}

	return nil
}

/*
itob transforms an int to a binary representation for BoltDB
*/
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
