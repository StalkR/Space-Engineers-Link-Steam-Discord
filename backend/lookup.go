package main

import (
  "fmt"
  "net/http"
  "strconv"
  "strings"
)

func init() {
  http.HandleFunc("/lookup/", lookupHandler)
}

func lookupHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path == "/lookup" {
    w.Header().Add("Content-Type", "text/plain")
    fmt.Fprintln(w, `GET /lookup/steam/123 returns discord ID`)
    fmt.Fprintln(w, `GET /lookup/discord/123 returns steam ID`)
    return
  }

  if !config.apiKeys[r.Header.Get("X-API-Key")] {
    http.Error(w, "you need a valid X-API-Key header", http.StatusForbidden)
    return
  }

  for path, f := range map[string]func(string) (string, bool){
    "/lookup/discord/": lookupByDiscordID,
    "/lookup/steam/":   lookupBySteamID,
  } {
    if !strings.HasPrefix(r.URL.Path, path) {
      continue
    }
    id := strings.TrimPrefix(r.URL.Path, path)
    if _, err := strconv.ParseUint(id, 10, 64); err != nil {
      http.Error(w, "invalid ID", http.StatusUnprocessableEntity)
      return
    }

    result, ok := f(id)
    if !ok {
      http.NotFound(w, r)
      return
    }
    w.Header().Add("Content-Type", "text/plain")
    fmt.Fprintln(w, result)
    return
  }

  http.NotFound(w, r)
}
