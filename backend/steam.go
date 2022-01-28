package main

import (
  "fmt"
  "net/http"
  "strconv"
  "strings"

  "github.com/StalkR/openid/openid20"
)

func init() {
  http.HandleFunc("/steam", steamHandler)
  http.HandleFunc("/steam/auth", steamAuthHandler)
  http.HandleFunc("/steam/check/", steamCheckHandler)
}

const steamEndpoint = "https://steamcommunity.com/openid/login"

func steamHandler(w http.ResponseWriter, r *http.Request) {
  target := openid20.RedirectURL(steamEndpoint, config.BaseURL+"/steam/auth?nonce="+setNonce(w))
  http.Redirect(w, r, target, http.StatusSeeOther)
}

func steamAuthHandler(w http.ResponseWriter, r *http.Request) {
  if err := verifyNonce(w, r, r.URL.Query().Get("nonce")); err != nil {
    http.Error(w, err.Error(), http.StatusForbidden)
    return
  }

  profileURL, err := openid20.Verify(r, steamEndpoint)
  if err != nil {
    http.Error(w, err.Error(), http.StatusForbidden)
    return
  }

  const prefix = "https://steamcommunity.com/openid/id/"
  if !strings.HasPrefix(profileURL, prefix) {
    http.Error(w, fmt.Sprintf("profile URL %v does not start with %v", profileURL, prefix), http.StatusInternalServerError)
    return
  }
  steamID := strings.TrimPrefix(profileURL, prefix)

  setSignedCookie(w, "steamID", steamID)
  http.Redirect(w, r, "/", http.StatusSeeOther)
}

func steamCheckHandler(w http.ResponseWriter, r *http.Request) {
  const path = "/steam/check/"
  if !strings.HasPrefix(r.URL.Path, path) {
    http.NotFound(w, r)
    return
  }
  id := strings.TrimPrefix(r.URL.Path, path)
  if _, err := strconv.ParseUint(id, 10, 64); err != nil {
    http.Error(w, "invalid ID", http.StatusUnprocessableEntity)
    return
  }

  if _, ok := lookupBySteamID(id); !ok {
    http.NotFound(w, r)
    return
  }
}
