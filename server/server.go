package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"raftybadger/badgerdb"
)

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func decode(text []byte) (KeyValuePair, error) {
	var pair KeyValuePair
	err := json.Unmarshal(text, pair)
	if err != nil {
		return KeyValuePair{}, err
	}
	return pair, nil
}

func getHandler(w http.ResponseWriter, r *http.Request, db *badgerdb.BadgerDB) {
	key, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 - Something bad happened!")
	}

	value, err := db.GetValue(string(key))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 - Could not find key %s", key)
	}

	fmt.Printf("Responding with value: %s\r\n", value)
	w.Write([]byte(value))
}

func setHandler(w http.ResponseWriter, r *http.Request, db *badgerdb.BadgerDB) {
	key, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 - Something bad happened!")
		return
	}

	var kvp KeyValuePair
	err = json.Unmarshal(key, &kvp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not decode json object: %s", key)
		return
	}

	err = db.SetValue(kvp.Key, kvp.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 - DB write failed with key: %s, value: %s", kvp.Key, kvp.Value)
	}

	w.WriteHeader(http.StatusOK)
}

func Serve(db *badgerdb.BadgerDB) {
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		getHandler(w, r, db)
	})
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		setHandler(w, r, db)
	})
	port := ":8090"
	fmt.Printf("Starting server on port: %s\r\n", port)
	http.ListenAndServe(port, nil)
}
