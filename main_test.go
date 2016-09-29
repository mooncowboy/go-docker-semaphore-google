package main

import (
  "testing"
  "time"
)

func TestCurrentTime(t *testing.T) {
  result := currentTime()
  // Check that err is null for RFC822Z time format
  _, err := time.Parse(time.RFC822Z, result)
  if err != nil {
    t.Fail()
  }
}
