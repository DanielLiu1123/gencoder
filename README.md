# Gencoder

The ultimate code generator for any languages/frameworks.

Applicable scenarios for gencoder:

- You need to modify the generated code and the modified code will not be overwritten

## Build

```bash
make && CGO_ENABLED=0 go build -o gencoder cmd/gencoder/main.go
```

If you updated Handlebars.js version or added new helpers in `pkg/jsruntime/helpers.js`:
```bash
go generate ./...
```