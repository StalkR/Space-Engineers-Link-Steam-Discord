package main

import (
  "context"
  "encoding/json"
  "fmt"
  "net/http"
  "net/url"

  "golang.org/x/oauth2"
)

func init() {
  http.HandleFunc("/discord", discordHandler)
  http.HandleFunc("/discord/auth", discordAuthHandler)
}

func discordHandler(w http.ResponseWriter, r *http.Request) {
  target := config.discordOAuth.AuthCodeURL(setNonce(w))
  http.Redirect(w, r, target, http.StatusSeeOther)
}

func discordAuthHandler(w http.ResponseWriter, r *http.Request) {
  ctx := context.Background()

  if err := verifyNonce(w, r, r.URL.Query().Get("state")); err != nil {
    http.Error(w, err.Error(), http.StatusForbidden)
    return
  }

  token, err := config.discordOAuth.Exchange(ctx, r.URL.Query().Get("code"))
  if err != nil {
    http.Error(w, "OAuth exchange error, please try again later", http.StatusForbidden)
    return
  }

  if err := discordAuth(ctx, w, r, token); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func discordAuth(ctx context.Context, w http.ResponseWriter, r *http.Request, token *oauth2.Token) error {
  client := config.discordOAuth.Client(ctx, token)
  resp, err := client.Get("https://discord.com/api/users/@me")
  if err != nil {
    return fmt.Errorf("users api: %v", err)
  }
  defer resp.Body.Close()
  if want := http.StatusOK; resp.StatusCode != want {
    return fmt.Errorf("users api: got status %v want %v", resp.Status, want)
  }
  var v struct {
    ID            string `json:"id"`
    Username      string `json:"username"`
    Discriminator string `json:"discriminator"`
  }
  if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
    return fmt.Errorf("json decode: %v", err)
  }

  // that's all we needed, we don't need to store the tokens
  // revoke to avoid appearing in "Authorized Apps"
  args := url.Values{}
  args.Set("client_id", config.Discord.ClientID)
  args.Set("client_secret", config.Discord.ClientSecret)
  args.Set("token", token.AccessToken)
  resp, err = http.PostForm("https://discord.com/api/oauth2/token/revoke", args)
  if err != nil {
    return err
  }
  if want := http.StatusOK; resp.StatusCode != want {
    return fmt.Errorf("revoke %v: got status %v want %v", "access token", resp.Status, want)
  }

  setSignedCookie(w, "discordID", v.ID)
  http.Redirect(w, r, "/", http.StatusFound)
  return nil
}
