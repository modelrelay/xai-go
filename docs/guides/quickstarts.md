# Service Quickstarts

## Embeddings

```go
resp, err := client.Embeddings.Embed(ctx, &xaiapiv1.EmbedRequest{
    Model: "embed-3",
    Input: []*xaiapiv1.EmbedInput{{Input: &xaiapiv1.EmbedInput_String{String_: "hello"}}},
})
if err != nil { log.Fatal(err) }
fmt.Println(resp.GetEmbeddings()[0].GetEmbeddings()[0].GetFloatArray())
```

## Images

```go
imgResp, err := client.Images.GenerateImage(ctx, &xaiapiv1.GenerateImageRequest{
    Model:  "grok-image-latest",
    Prompt: "A rocket made of code",
})
if err != nil { log.Fatal(err) }
fmt.Println(imgResp.GetImages()[0].GetUrl())
```

## Documents

```go
searchResp, err := client.Documents.Search(ctx, documents.SearchRequest(
    "onboarding guide",
    documents.CollectionSource("engineering"),
    documents.WithLimit(5),
))
if err != nil { log.Fatal(err) }
for _, match := range searchResp.GetMatches() {
    fmt.Println(match.GetChunkContent())
}
```

## Tokenize

```go
tokResp, err := client.Tokenize.TokenizeText(ctx, &xaiapiv1.TokenizeTextRequest{
    Model: "grok-2-latest",
    Text:  "Hello world",
})
if err != nil { log.Fatal(err) }
fmt.Println(tokResp.GetTokens())
```

## Models

```go
modelsResp, err := client.Models.ListLanguageModels(ctx)
if err != nil { log.Fatal(err) }
for _, model := range modelsResp.GetModels() {
    fmt.Println(model.GetName())
}
```

## Auth

```go
keyInfo, err := client.Auth.GetAPIKeyInfo(ctx)
if err != nil { log.Fatal(err) }
fmt.Println(keyInfo.GetName(), keyInfo.GetTeamId())
```
