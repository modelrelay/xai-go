package encrypted

import (
	"fmt"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// Decoder converts encrypted payloads to plaintext.
type Decoder func(ciphertext string) (string, error)

// DecryptResponse walks a GetChatCompletionResponse and replaces encrypted content using the decoder.
func DecryptResponse(resp *xaiapiv1.GetChatCompletionResponse, decode Decoder) error {
	if resp == nil || decode == nil {
		return nil
	}
	for _, output := range resp.GetOutputs() {
		msg := output.GetMessage()
		if msg == nil {
			continue
		}
		if err := applyToMessage(msg, decode); err != nil {
			return err
		}
	}
	return nil
}

// DecryptChunk applies the decoder to encrypted content in a streaming chunk delta.
func DecryptChunk(chunk *xaiapiv1.GetChatCompletionChunk, decode Decoder) error {
	if chunk == nil || decode == nil {
		return nil
	}
	for _, out := range chunk.GetOutputs() {
		delta := out.GetDelta()
		if delta == nil {
			continue
		}
		if delta.GetEncryptedContent() != "" {
			plain, err := decode(delta.GetEncryptedContent())
			if err != nil {
				return fmt.Errorf("decrypt chunk delta: %w", err)
			}
			delta.EncryptedContent = ""
			delta.Content += plain
		}
		for _, call := range delta.GetToolCalls() {
			if err := applyToTool(call, decode); err != nil {
				return err
			}
		}
	}
	return nil
}

func applyToMessage(msg *xaiapiv1.CompletionMessage, decode Decoder) error {
	if msg.GetEncryptedContent() == "" {
		return nil
	}
	plain, err := decode(msg.GetEncryptedContent())
	if err != nil {
		return fmt.Errorf("decrypt message: %w", err)
	}
	msg.EncryptedContent = ""
	msg.Content += plain
	for _, call := range msg.GetToolCalls() {
		if err := applyToTool(call, decode); err != nil {
			return err
		}
	}
	return nil
}

func applyToTool(call *xaiapiv1.ToolCall, decode Decoder) error {
	if call == nil {
		return nil
	}
	fn := call.GetFunction()
	if fn == nil || fn.GetArguments() == "" {
		return nil
	}
	plain, err := decode(fn.GetArguments())
	if err != nil {
		return fmt.Errorf("decrypt tool args: %w", err)
	}
	fn.Arguments = plain
	return nil
}
