---
title: "How to hash and compare passwords in Go"
publishdate: 2023-12-23
categories: 
    - Programming
tags:
  - go
  - bcrypt
  - security
---

The best to hash passwords in Go is using `golang.org/x/crypto/bcrypt`:

```go
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
  
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

You should use the default `bcrypt.DefaultCost` just in case that the current value will become not sufficient and the default cost will increase.
