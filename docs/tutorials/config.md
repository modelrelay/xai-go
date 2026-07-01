# Declarative Configuration

Use `config.Config` when you want to load settings (API key, base URL, user agent) from JSON/YAML instead of threading option functions through your application.

```go
type Settings struct {
    Grok config.Config `json:"grok"`
}

var s Settings
if err := json.Unmarshal(data, &s); err != nil {
    log.Fatal(err)
}

client, err := s.Grok.NewClient(context.Background())
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

You can still override individual options by passing additional `xai.Option` values to `NewClient`.
