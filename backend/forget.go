package main

import (
  "net/http"
)

func init() {
  http.HandleFunc("/forget", forgetHandler)
}

func forgetHandler(w http.ResponseWriter, r *http.Request) {
  steamID, err := getSignedCookie(r, "steamID")
  if err != nil {
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
  }
  if err := removeMapping(steamID); err != nil {
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
  }
  deleteCookie(w, "steamID")
  deleteCookie(w, "discordID")
  http.Redirect(w, r, "/", http.StatusSeeOther)
}
