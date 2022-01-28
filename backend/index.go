package main

import (
  "embed"
  "log"
  "net/http"
  "text/template"
)

var (
  //go:embed img
  img embed.FS
  //go:embed templates
  templates embed.FS
)

func init() {
  http.HandleFunc("/", indexHandler)
  http.Handle("/img/", http.FileServer(http.FS(img)))
}

var indexTmpl = template.Must(template.ParseFS(templates, "templates/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }
  steamID, err := getSignedCookie(r, "steamID")
  if err != nil {
    steamID = ""
  }
  discordID, err := getSignedCookie(r, "discordID")
  if err != nil {
    discordID = ""
  }

  if steamID != "" && discordID != "" {
    if err := saveMapping(steamID, discordID); err != nil {
      log.Printf("save mapping: %v", err)
      http.Error(w, "failed to save mapping", http.StatusInternalServerError)
      return
    }
  }

  w.Header().Add("Content-Type", "text/html")
  if err := indexTmpl.Execute(w, struct {
    SteamID, DiscordID string
  }{
    SteamID:   steamID,
    DiscordID: discordID,
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}
