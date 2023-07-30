# MEASUREYOURLIFE Web Service

## Build, run and test

Run unit-tests

```
make test
```

Run integration tests (The app needs to be running!)

```
make integration_test
```

Run the project

```
make run
```

Run project with live updates while developing

```
gowatch
```

## Environment Variables

```
MEASUREYOURLIFE_PORT=:8700
MEASUREYOURLIFE_ALLOW_ORIGIN=http://127.0.0.1:8080
MEASUREYOURLIFE_SESSION_ENCRYPTION_PASSPHRASE=<some secret passphrase>

MEASUREYOURLIFE_TLS=false
MEASUREYOURLIFE_CERT_FILE=cert.pem
MEASUREYOURLIFE_KEY_FILE=key.unencrypted.pem
```

## API
