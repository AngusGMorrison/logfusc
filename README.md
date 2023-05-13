# logfusc

Instrument your codebase with confidence. logfusc is a Go library that makes
redacting sensitive data from logs and traces simple. Stop scrubbing logs in the
aftermath of preventable human errors. Make your secrets **unloggable** instead.

## logfusc.Secret

`logfusc.Secret` is a generic wrapper for any type that you want to redact from
logs, traces and other outputs.

`Secret` implements `fmt.Stringer` and `fmt.GoStringer`, so no matter how hard
you try to format it, it doesn't give up its secret.

```go
password := "do not log!"
secret := logfusc.Secret(password)

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

### Log with reckless abandon

`logfusc.Secret` redacts your secrets when marshaled to a variety of formats, so
you can pass complete structs to your logger without worrying about leaking
sensitive data. No more manual redaction. No configuration. No leaks.

```go
type User struct {
    Name     string `json:"name"`
    SecretOfLife logfusc.Secret[int] `json:"secret_of_life"`
}

user := User{
    Name:     "alice",
    SecretOfLife: logfusc.NewSecret(42),
}

b, _ := json.Marshal(user)

log.Println(string(b))

// => {"name":"alice","secret_of_life":"logfusc.Secret[int]{REDACTED}"}
```

So far, `Secret` satisfies:
- `json.Marshaler`
- More coming soon! Too slow for you? Why not contribute?

### Decode directly to logfusc.Secret

Protect secrets at the boundaries of your service by decoding them directly into `logfusc.Secret`.

So far, `Secret` satisfies:
- `json.Unmarshaler`
- More coming soon! Too slow for you? Why not contribute?

### Use secrets with intention

Secret encourages you to log often and trace **everything** in the knowledge
that your secrets are safe. That doesn't stop you working with sensitive data,
though, but you have to do it *deliberately*.

```go
secret := logfusc.NewSecret("PEM string")
log.Println(secret.Expose())

// => PEM string
```

## Compatibility

`logfusc.Secret` is tested and known to redact its contents from the following
loggers:
- `log`
- More compatibility tests coming soon! Too slow for you? Why not contribute?

## Caveats

`logfusc` won't protect your secrets from third-party code that really digs for
them. In Go, you can access the private fields of a struct using a combination
of the `reflect` and `unsafe` packages. But if your telemetry package is doing
that, then you should probably find another one...
