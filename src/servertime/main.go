package main

import (
    "time"
    "encoding/json"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  t := time.Now()
  js, err := json.Marshal(t)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
