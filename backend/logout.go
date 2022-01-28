package main

import (
  "net/http"
)

func init() {
  http.HandleFunc("/logout", logoutHandler)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
  deleteCookie(w, "steamID")
  deleteCookie(w, "discordID")
  http.Redirect(w, r, "/", http.StatusSeeOther)
}
