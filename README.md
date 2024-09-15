# Gencoder

## Build

```bash
make && CGO_ENABLED=0 go build -o gencoder cmd/gencoder/main.go
```

If you want to update Handlebars.js version or add new helpers in `pkg/jsruntime/helpers.js`:
```bash
make hb
```