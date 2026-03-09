package main

import (
  "fmt"
  "time"

  jwt "github.com/golang-jwt/jwt/v5"
)

func main() {
  uid := "69ade7e93429c14f71db5b4b"
  login := "qa1773004777"
  secret := "dev-jwt-secret-local-only"

  now := time.Now().UTC()
  claims := jwt.MapClaims{
    "uid": uid,
    "login": login,
    "role": "admin",
    "sub": uid,
    "iat": now.Unix(),
    "exp": now.Add(24 * time.Hour).Unix(),
  }

  t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  s, err := t.SignedString([]byte(secret))
  if err != nil {
    panic(err)
  }

  fmt.Print(s)
}
