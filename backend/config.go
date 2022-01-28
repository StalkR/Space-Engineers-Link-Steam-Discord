package main

import (
	"encoding/json"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

// config contains the loaded config file.
// It is in a file not flags because it contains secrets.
var config struct {
	// Address and port to listen on.
	// Example: ":8080"
	Listen string `json:"listen"`

	// Application base URL to configure redirect endpoints.
	// Example: "https://example.com"
	BaseURL string `json:"base_url"`

	// Application secret to sign and verify state. 16 bytes of random is good.
	// Example: "eiQuoo[la:uB`ae6"
	Secret string `json:"secret"`

	// Discord OAuth config.
	// Example: {"client_id": "123456", "client_secret": "xxxxx"}
	Discord struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"discord"`
	discordOAuth *oauth2.Config `json:"-"` // once loaded

	// API keys for /lookup access. 16 random alphanum is good.
	// Example: ["shaeT4eel7yoh8ka", "baemahpai0Yae6ao"]
	APIKeys []string        `json:"api_keys"`
	apiKeys map[string]bool `json:"-"` // once loaded

	// SQLite3 database file where to store the "steamID,discordID" mappings.
	// Example: "mappings.db"
	Database string `json:"database"`
}

func loadConfig(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return err
	}

	config.BaseURL = strings.TrimSuffix(config.BaseURL, "/") // just in case

	config.discordOAuth = &oauth2.Config{
		ClientID:     config.Discord.ClientID,
		ClientSecret: config.Discord.ClientSecret,
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		RedirectURL: config.BaseURL + "/discord/auth",
	}

	config.apiKeys = map[string]bool{}
	for _, v := range config.APIKeys {
		config.apiKeys[v] = true
	}

	return nil
}
