package main

import (
  "time"
  "encoding/json"
  "net/http"
)

type ServiceResult struct {
  FormattedTime string
  Greeting string
}

func currentTime() string {
  return time.Now().Format(time.RFC822Z)
}

func handler(w http.ResponseWriter, r *http.Request) {
  t := currentTime()
  sr := ServiceResult{t, "Hi there"}
  js, err := json.Marshal(sr)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func main() {
  http.HandleFunc("/", handler)
  http.ListenAndServe(":80", nil)
}
