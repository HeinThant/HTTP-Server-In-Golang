package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	items  = make(map[int]Item)
	nextID = 1
	mutex  = sync.Mutex{}
)

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	defer mutex.Unlock()

	var itemList []Item
	for _, item := range items {
		itemList = append(itemList, item)
	}
	json.NewEncoder(w).Encode(itemList)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	newItem.ID = nextID
	items[nextID] = newItem
	nextID++
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	if _, exists := items[id]; !exists {
		mutex.Unlock()
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	updatedItem.ID = id
	items[id] = updatedItem
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	if _, exists := items[id]; !exists {
		mutex.Unlock()
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	delete(items, id)
	mutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func resetServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	mutex.Lock()
	items = make(map[int]Item) // Clear the map
	nextID = 1                 // Reset ID counter
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server reset successfully"))
}

func main() {
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getItems(w, r)
		case "POST":
			createItem(w, r)
		case "PUT":
			updateItem(w, r)
		case "DELETE":
			deleteItem(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/reset", resetServer) // Add new reset route
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.ListenAndServe(":8080", nil)
}
