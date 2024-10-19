package main

// Website, recording viewers that go to the website 
// Everyday we write to the log file (ts pager id and customer id)
// Each day will be a different log file 
// Generates a list of customers that meet a criteria
// showed up on two days, they visited at least two unique pages 
// 2 diff file (two slices: day 1 and day 2)

import (
  "time"
  "net/http"
  "github.com/gorilla/mux"
  "log"
  "encoding/json"
  "os"
  "sync"
  )

type View struct {
  mu sync.Mutex
  CustomerID string `json:"customer_id"`
  PageID string `json:"page_id"`
  Timestamp time.Time `json:"timestamp"`
}

type Views struct {
  Day1 []View  `json:"day1"`
  Day2 []View `json:"day2"`
}

type CustomerID string


func loadViews(filename string) (Views, error) {
    var views Views
    data, err := os.ReadFile(filename)
    if err != nil {
        return views, err
    }
    
    err = json.Unmarshal(data, &views)
    return views, err
}

func logView(w http.ResponseWriter, r *http.Request) {
	var view View
	err := json.NewDecoder(r.Body).Decode(&view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	view.Timestamp = time.Now()

	view.mu.Lock()
	defer view.mu.Unlock()

	// // Determine which day's slice to append to
	// if time.Now().Day()%2 == 0 {
	// 	view.Day2 = append(view.Day2, view)
	// 	writeToLogFile("day2.log", view)
	// } else {
	// 	view.Day1 = append(view.Day1, view)
	// 	writeToLogFile("day1.log", view)
	// }

	w.WriteHeader(http.StatusCreated)
}

func main() {

  // Initialize and start the gateway
  router := mux.NewRouter()

  // Add your routes here, for example:
  	// Add routes
	router.HandleFunc("/view", logView).Methods("POST")
	// router.HandleFunc("/analyze", analyzeViews).Methods("GET")

  log.Fatal(http.ListenAndServe(":8080", router))
}