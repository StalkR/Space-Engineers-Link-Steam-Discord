package main

import (
  "database/sql"
  "fmt"
  "log"
  "strconv"

  _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func loadMappings() error {
  var err error
  db, err = sql.Open("sqlite3", config.Database)
  if err != nil {
    return err
  }
  if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS mappings(
          steamID INTEGER PRIMARY KEY,
          discordID INTEGER NOT NULL
      );`); err != nil {
    return err
  }
  var total int
  if err := db.QueryRow(`SELECT COUNT(*) FROM mappings`).Scan(&total); err != nil {
    return err
  }
  log.Printf("Database mappings: %v", total)
  return nil
}

func saveMapping(steamID, discordID string) error {
  steamIDu, err := strconv.ParseUint(steamID, 10, 64)
  if err != nil {
    return err
  }
  discordIDu, err := strconv.ParseUint(discordID, 10, 64)
  if err != nil {
    return err
  }
  _, err = db.Exec(`REPLACE INTO mappings (steamID, discordID) VALUES (?, ?)`, steamIDu, discordIDu)
  return err
}

func removeMapping(steamID string) error {
  steamIDu, err := strconv.ParseUint(steamID, 10, 64)
  if err != nil {
    return err
  }
  _, err = db.Exec(`DELETE FROM mappings WHERE steamID = ?`, steamIDu)
  return err
}

func lookupBySteamID(id string) (string, bool) {
  n, err := strconv.ParseUint(id, 10, 64)
  if err != nil {
    return "", false
  }
  row := db.QueryRow(`SELECT discordID FROM mappings WHERE steamID = ?`, n)
  var r uint64
  if err := row.Scan(&r); err != nil {
    return "", false
  }
  return fmt.Sprintf("%v", r), true
}

func lookupByDiscordID(id string) (string, bool) {
  n, err := strconv.ParseUint(id, 10, 64)
  if err != nil {
    return "", false
  }
  row := db.QueryRow(`SELECT steamID FROM mappings WHERE discordID = ?`, n)
  var r uint64
  if err := row.Scan(&r); err != nil {
    return "", false
  }
  return fmt.Sprintf("%v", r), true
}
