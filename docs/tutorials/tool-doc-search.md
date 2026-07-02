# Tool + Document Search Walkthrough

This short walkthrough shows how to let Grok call a custom `lookup_docs` tool, run a document search, and feed the results back to the model.

1. **Prepare a tool definition** using the helpers in `tools`:

   ```go
   fnTool, err := tools.FunctionTool("lookup_docs", "Look up documents", map[string]any{
       "type": "object",
       "properties": map[string]any{
           "query": map[string]any{"type": "string"},
       },
       "required": []string{"query"},
   })
   if err != nil {
       log.Fatal(err)
   }
   ```

2. **Kick off a streaming request** with the tool attached:

   ```go
   stream, err := client.Responses.CreateStream(ctx, &xaiapiv1.GetCompletionsRequest{
       Model: "grok-4.3",
       Messages: []*xaiapiv1.Message{
           messages.SystemText("Use lookup_docs when you need internal knowledge."),
           messages.UserText("Summarize the engineering onboarding guide."),
       },
       Tools: []*xaiapiv1.Tool{fnTool},
   })
   if err != nil {
       log.Fatal(err)
   }
   tracker := responses.NewToolCallTracker()
   registry := toolruntime.NewRegistry()
   ```

3. **Register a handler** that runs document search when the tool is invoked:

   ```go
   registry.Register("lookup_docs", func(ctx context.Context, fn *xaiapiv1.FunctionCall) (any, error) {
       resp, err := client.Documents.Search(ctx, documents.SearchRequest(fn.GetArguments(), documents.CollectionSource("onboarding")))
       if err != nil {
           return nil, err
       }
       return resp.GetMatches(), nil
   })
   ```

4. **Stream chunks, detect tool events, and feed ROLE_TOOL messages back**:

   ```go
   acc := responses.NewAccumulator()
   if err := stream.ForEachChunk(ctx, func(chunk *xaiapiv1.GetChatCompletionChunk) error {
       acc.AddChunk(chunk)
       for _, event := range tracker.ConsumeChunk(chunk) {
           if !event.Complete {
               continue
           }
           msg, err := registry.Handle(ctx, event)
           if err != nil {
               return err
           }
           _ = msg // Append msg to your conversation state before resuming the loop.
       }
       return nil
   }); err != nil {
       log.Fatal(err)
   }
   ```

5. **When streaming finishes**, inspect the accumulator for the final reply:

   ```go
   outs := acc.Response().GetOutputs()
   if len(outs) == 0 {
       log.Fatal("stream produced no output")
   }
   fmt.Println(outs[0].GetMessage().GetContent())
   ```

The full runnable version of this example lives in `examples/tool_call`. Adjust the collection IDs and tool logic to fit your corpus.
