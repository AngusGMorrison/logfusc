# logfusc – surefire secret redaction for logs and traces

[![Go
Reference](https://pkg.go.dev/badge/github.com/angusgmorrison/logfusc.svg)](https://pkg.go.dev/github.com/angusgmorrison/logfusc) ![CI](https://github.com/AngusGMorrison/logfusc/actions/workflows/vet.yml/badge.svg)

Instrument your codebase with confidence. logfusc is a Go library that makes
redacting sensitive data from logs and traces simple. Stop scrubbing logs in the
aftermath of preventable human errors. Make your secrets **unloggable**.

![A family of gophers riding a log down a river](images/logging-gophers.jpg)


## logfusc.Secret

`logfusc.Secret` is a generic wrapper for any type that you want to redact from
logs, traces and other outputs.

`Secret` implements `fmt.Stringer` and `fmt.GoStringer`, so no matter how hard
you try to format it, it doesn't give up its secret.

```go
password := "do not log!"
secret := logfusc.NewSecret(password)

fmt.Printf("%s\n", secret)
// => logfusc.Secret[string]{REDACTED}

fmt.Printf("%q\n", secret)
// => "logfusc.Secret[string]{REDACTED}"

fmt.Printf("%v\n", secret)
// => "logfusc.Secret[string]{REDACTED}"

fmt.Printf("%+v\n", secret)
// => "logfusc.Secret[string]{REDACTED}"

fmt.Printf("%#v\n", secret)
// => "logfusc.Secret[string]{REDACTED}"

fmt.Printf("%x\n", secret)
// => 6c6f67667573632e5365637265745b737472696e675d7b52454441435445447d == logfusc.Secret[string]{REDACTED}
```

### Log anything, anywhere

`logfusc.Secret` redacts your secrets when marshaled to a variety of formats, so
you can pass complete structs to your logger without worrying about leaking
sensitive data. No more manual redaction. No configuration. No leaks.

```go
type Universe struct {
    SecretOfLife logfusc.Secret[int] `json:"secret_of_life"`
}

func main() {
  universe := Universe{
      SecretOfLife: logfusc.NewSecret(42),
  }

  b, _ := json.Marshal(universe)

  log.Println(string(b))
}

// => {"name":"alice","secret_of_life":"logfusc.Secret[int]{REDACTED}"}
```

So far, `Secret` satisfies:
- `json.Marshaler` (tested with both `encoding/json` and [json-iterator](https://github.com/json-iterator/go))
- More coming soon! Too slow for you? Why not contribute?

### Decode directly to logfusc.Secret

Protect secrets at the boundaries of your service by decoding them directly into
`logfusc.Secret`.

```go
type RegisterRequest struct {
  Email string `json:"email"`
  Password logfusc.Secret[string] `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
  var req RegisterRequest
  _ = json.NewDecoder(r.Body).Decode(&req)
  fmt.Println(req.Password)
  fmt.Println(req.Password.Expose())
}

// => logfusc.Secret[string]{REDACTED}
//    password
```

So far, `Secret` satisfies:
- `json.Unmarshaler`
- More coming soon! Too slow for you? Why not contribute?

### Use secrets with intention

`Secret` encourages you to log often and trace **everything** in the knowledge that
your secrets are safe. That doesn't stop you working with sensitive data,
but you have to do it *deliberately*.

```go
secret := logfusc.NewSecret("PEM string")
log.Println(secret.Expose())

// => PEM string
```

## Compatibility

When logged as a standalone object or as an exported struct field,
`logfusc.Secret` will work as described with all sane logging libraries.

For your peace of mind, `logfusc.Secret` has been explicitly tested for
compatibility with the following loggers:
- `log`
- `logrus`
- `zap`
- `zerolog`
- More compatibility tests coming soon! Too slow for you? Why not contribute?

### ⚠️ Caution ⚠️
The standard library `fmt` package and any third-party packages that use reflection
to deeply print structs __DO NOT__ respect `fmt.Stringer` or `fmt.GoStringer`
implementations for unexported struct fields.

If a `logfusc.Secret` is logged as part of a struct with unexported fields, it
__will__ be logged in plain text.

Alternatively:
* Use `logfusc.Secret` only in exported struct fields.
* Define custom `fmt.Stringer` and `fmt.GoStringer` implementations for structs
  that contain unexported `logfusc.Secret` fields, invoking `Secret.String()` or
  `Secret.GoString()` manually for these fields.
* Use [go-spew](https://github.com/davecgh/go-spew) instead of `fmt`, which deeply prints structs without reflection,
  respecting `fmt.Stringer` and `fmt.GoStringer` implementations for unexported
  fields.
