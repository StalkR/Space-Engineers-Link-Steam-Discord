package main

import (
  "crypto/hmac"
  "crypto/rand"
  "crypto/sha256"
  "encoding/hex"
  "fmt"
  "net/http"
)

func getCookie(r *http.Request, name string) (string, error) {
  c, err := r.Cookie("__Host-" + name)
  if err != nil {
    return "", err
  }
  return c.Value, nil
}

func setCookieWithAge(w http.ResponseWriter, name, value string, maxAge int) {
  http.SetCookie(w, &http.Cookie{
    Name:  "__Host-" + name,
    Value: value,
    Path:  "/",
    // Secure determines if cookies have the secure flag, requiring https.
    // If you run the app over http and it doesn't work, that's why.
    Secure:   true,
    HttpOnly: true,
    MaxAge:   maxAge,
  })
}

func setCookie(w http.ResponseWriter, name, value string) {
  setCookieWithAge(w, name, value, 0)
}

func deleteCookie(w http.ResponseWriter, name string) {
  setCookieWithAge(w, name, "", -1)
}

func getSignedCookie(r *http.Request, name string) (string, error) {
  c, err := getCookie(r, name)
  if err != nil {
    return "", err
  }
  return verify(c)
}

func setSignedCookie(w http.ResponseWriter, name, value string) {
  setCookie(w, name, sign(value))
}

func randBytes(n int) []byte {
  b := make([]byte, n)
  if _, err := rand.Read(b); err != nil {
    panic(err)
  }
  return b
}

const nonceLength = 16

func sign(m string) string {
  nonce := randBytes(nonceLength)
  mac := hmac.New(sha256.New, []byte(config.Secret))
  mac.Write([]byte(m))
  mac.Write(nonce)
  return hex.EncodeToString(mac.Sum(nil)) + hex.EncodeToString(nonce) + m
}

func verify(m string) (string, error) {
  if len(m) < sha256.Size*2+nonceLength*2 {
    return "", fmt.Errorf("message too short")
  }
  expectedMAC, err := hex.DecodeString(m[:sha256.Size*2])
  if err != nil {
    return "", err
  }
  m = m[sha256.Size*2:]
  nonce, err := hex.DecodeString(m[:nonceLength*2])
  if err != nil {
    return "", err
  }
  m = m[nonceLength*2:]

  mac := hmac.New(sha256.New, []byte(config.Secret))
  mac.Write([]byte(m))
  mac.Write(nonce)
  if !hmac.Equal(mac.Sum(nil), expectedMAC) {
    return "", fmt.Errorf("invalid MAC")
  }
  return m, nil
}

func setNonce(w http.ResponseWriter) string {
  nonce := hex.EncodeToString(randBytes(nonceLength))
  setCookie(w, "nonce", nonce)
  return nonce
}

func verifyNonce(w http.ResponseWriter, r *http.Request, got string) error {
  want, err := getCookie(r, "nonce")
  if err != nil {
    return fmt.Errorf("missing nonce")
  }
  deleteCookie(w, "nonce")
  if got != want {
    return fmt.Errorf("invalid nonce")
  }
  return nil
}
