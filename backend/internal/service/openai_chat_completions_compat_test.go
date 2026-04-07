package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertChatCompletionsRequest(t *testing.T) {
	reqBody := map[string]any{
		"model": "gpt-5.4",
		"messages": []any{
			map[string]any{
				"role":    "system",
				"content": "Be concise.",
			},
			map[string]any{
				"role":    "user",
				"content": "hello",
			},
			map[string]any{
				"role":    "assistant",
				"content": "working",
				"tool_calls": []any{
					map[string]any{
						"id":   "call_1",
						"type": "function",
						"function": map[string]any{
							"name":      "weather",
							"arguments": "{\"city\":\"Shanghai\"}",
						},
					},
				},
			},
			map[string]any{
				"role":         "tool",
				"tool_call_id": "call_1",
				"content":      "sunny",
			},
		},
		"tools": []any{
			map[string]any{
				"type": "function",
				"function": map[string]any{
					"name":        "weather",
					"description": "Get weather",
					"parameters":  map[string]any{"type": "object"},
				},
			},
		},
		"tool_choice": map[string]any{
			"type": "function",
			"function": map[string]any{
				"name": "weather",
			},
		},
		"max_tokens": 128,
	}

	modified, err := ConvertChatCompletionsRequest(reqBody)
	require.NoError(t, err)
	require.True(t, modified)

	require.NotContains(t, reqBody, "messages")
	require.Equal(t, "Be concise.", reqBody["instructions"])
	require.Equal(t, 128, reqBody["max_output_tokens"])

	tools, ok := reqBody["tools"].([]any)
	require.True(t, ok)
	require.Len(t, tools, 1)
	tool := tools[0].(map[string]any)
	require.Equal(t, "weather", tool["name"])
	require.Equal(t, "Get weather", tool["description"])

	toolChoice, ok := reqBody["tool_choice"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "function", toolChoice["type"])
	require.Equal(t, "weather", toolChoice["name"])

	input, ok := reqBody["input"].([]any)
	require.True(t, ok)
	require.Len(t, input, 4)

	userMsg := input[0].(map[string]any)
	require.Equal(t, "user", userMsg["role"])
	require.Equal(t, "hello", userMsg["content"])

	assistantMsg := input[1].(map[string]any)
	require.Equal(t, "assistant", assistantMsg["role"])
	require.Equal(t, "working", assistantMsg["content"])

	functionCall := input[2].(map[string]any)
	require.Equal(t, "function_call", functionCall["type"])
	require.Equal(t, "call_1", functionCall["call_id"])
	require.Equal(t, "weather", functionCall["name"])

	functionOutput := input[3].(map[string]any)
	require.Equal(t, "function_call_output", functionOutput["type"])
	require.Equal(t, "call_1", functionOutput["call_id"])
	require.Equal(t, "sunny", functionOutput["output"])
}

func TestConvertResponsesBodyToChatCompletions(t *testing.T) {
	service := &OpenAIGatewayService{
		toolCorrector: NewCodexToolCorrector(),
	}

	body := []byte(`{
		"id": "resp_123",
		"model": "gpt-5.4",
		"created_at": "2026-04-08T01:48:00Z",
		"output": [
			{
				"type": "message",
				"role": "assistant",
				"content": [
					{"type": "output_text", "text": "API is working"}
				]
			}
		],
		"usage": {
			"input_tokens": 11,
			"output_tokens": 7,
			"input_tokens_details": {
				"cached_tokens": 3
			}
		}
	}`)

	converted, usage, err := service.convertResponsesBodyToChatCompletions(body, "gpt-5.4")
	require.NoError(t, err)
	require.Equal(t, 11, usage.InputTokens)
	require.Equal(t, 7, usage.OutputTokens)
	require.Equal(t, 3, usage.CacheReadInputTokens)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(converted, &payload))
	require.Equal(t, "resp_123", payload["id"])
	require.Equal(t, "chat.completion", payload["object"])
	require.Equal(t, "gpt-5.4", payload["model"])

	choices := payload["choices"].([]any)
	require.Len(t, choices, 1)
	choice := choices[0].(map[string]any)
	require.Equal(t, "stop", choice["finish_reason"])
	message := choice["message"].(map[string]any)
	require.Equal(t, "assistant", message["role"])
	require.Equal(t, "API is working", message["content"])

	usagePayload := payload["usage"].(map[string]any)
	require.Equal(t, float64(11), usagePayload["prompt_tokens"])
	require.Equal(t, float64(7), usagePayload["completion_tokens"])
	require.Equal(t, float64(18), usagePayload["total_tokens"])
}

func TestConvertResponsesStreamEventToChatCompletions(t *testing.T) {
	service := &OpenAIGatewayService{
		toolCorrector: NewCodexToolCorrector(),
	}
	state := newChatCompletionsStreamState("chatcmpl-test", "gpt-5.4")

	lines, err := service.convertResponsesStreamEventToChatCompletions(`{"type":"response.output_text.delta","delta":"Hello"}`, state)
	require.NoError(t, err)
	require.Len(t, lines, 2)
	require.Contains(t, lines[0], `"role":"assistant"`)
	require.Contains(t, lines[1], `"content":"Hello"`)

	lines, err = service.convertResponsesStreamEventToChatCompletions(`{"type":"response.completed","response":{"id":"resp_done","model":"gpt-5.4"}}`, state)
	require.NoError(t, err)
	require.Len(t, lines, 2)
	require.Contains(t, lines[0], `"finish_reason":"stop"`)
	require.Equal(t, "data: [DONE]\n\n", lines[1])
}
