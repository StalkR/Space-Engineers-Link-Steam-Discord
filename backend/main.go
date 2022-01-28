package main

import (
  "flag"
  "log"
  "net/http"
)

var (
  flagConfig = flag.String("config", "config.json", "App configuration.")
)

func main() {
  flag.Parse()
  if err := loadConfig(*flagConfig); err != nil {
    log.Fatal(err)
  }
  if err := loadMappings(); err != nil {
    log.Fatal(err)
  }
  log.Printf("Listening on %v", config.Listen)
  if err := http.ListenAndServe(config.Listen, nil); err != nil {
    log.Fatal(err)
  }
}
