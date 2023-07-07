---
title: "Password policies"
publishdate: 2023-07-07
categories:
  - Ecommerce
tags:
  - auth
  - golang
---

One of tasks I was working on recently is related to [password policies](https://github.com/golang-app/ecommerce/issues/44). Of course, everything is configurable in code right now. In this note I want to tell you about some my decisions and how I got to them.

My first idea was creating an interface that any policy will have to satisfy.

```go
type PasswordPolicy interface {
    Verify(pass string) error
}
```

That make sense, doesn't it? When I was working on specific policy implementation I had a feeling that the type doesn't have to be an interface. A regular function should be enough. I rewroted it into fuctions then.

```go
type PasswordPolicy func(string) error
```

What if we need to make them configurable? Nothing simpler, let's use more functional-style code.

```go
func MinLength(n int) PasswordPolicy {
	return func(password string) error {
		if len(password) < n {
			return ErrPasswordTooShort
		}
		return nil
	}
}
```

Another thing that I changed a few times are errors. At the very beginning, I had only one error declared

```go
var ErrPasswordTooWeak = errors.New("passowrd is too weak")
```

This error code worked OK for very long but when I added a check for maximal password length the error didn't fit. `password is too weak` sounds very generic too. The user wouldn't know **why** the code is considered weak. I started adding more specific errors vor each function as shown below.

```go
var ErrPasswordTooShort = errors.New("password is too short")
var ErrPasswordLeaked = errors.New("password leaked")
var ErrPasswordTooLong = errors.New("password too long")
var ErrPasswordDoesNotContainLowercase = errors.New("password does not contain lowercase letter")
var ErrPasswordDoesNotContainUppercase = errors.New("password does not contain uppercase letter")
var ErrPasswordDoesNotContainNumber = errors.New("password does not contain number")
var ErrPasswordDoesNotContainSpecialChar = errors.New("password does not contain special character")
```

The good point is that every policy will have a specific error so it will be easy to test it. On the other hand, the number of errors is increasing what doesn't have a good impact on the readability. After some time, I decided to replace all errors with a custom type.


```go
type PasswordPolicyError string

func (e PasswordPolicyError) Error() string {
	return string(e)
}
```

It implements the `error` interface so it can be used as a regular error. In every place where I want to return an error, I just create a new instance of the custom type with a proper error message.

```go
func MinLength(n int) PasswordPolicy {
	return func(password string) error {
		if len(password) < n {
			return PasswordPolicyError("password is too short")
		}
		return nil
	}
}
```

Probably, I'll have to extend this type but until now it looks fine for me.

I have plans to add another policy that will restrict using a password that [already leaked](https://github.com/golang-app/ecommerce/issues/48) but I have to postpone the idea until I resolve an issue with access to the API. The [haveibeenpwned](https://haveibeenpwned.com/API/Key) API is quite expensive. I have plans to write a small service where I'll download all leaked passwords and expose them as a free API. The database is quite big (about 200 gb) so I have to import it in parts.
