package service

import (
	"encoding/json"
	"strings"
)

// ClaudeToOpenAIRequest 将 Claude Messages API 请求转换为 OpenAI Responses API 请求
// Claude: {"model":"claude-xxx","messages":[...],"system":"...","max_tokens":1024}
// OpenAI: {"model":"gpt-xxx","input":[...],"instructions":"..."}
func ClaudeToOpenAIRequest(claudeBody []byte) ([]byte, error) {
	var claudeReq map[string]any
	if err := json.Unmarshal(claudeBody, &claudeReq); err != nil {
		return nil, err
	}

	openaiReq := make(map[string]any)

	// 转换 model
	if model, ok := claudeReq["model"].(string); ok {
		openaiReq["model"] = normalizeCodexModel(model)
	}

	// 转换 system -> instructions
	if system, ok := claudeReq["system"].(string); ok && system != "" {
		openaiReq["instructions"] = system
	}

	// 转换 messages -> input
	if messages, ok := claudeReq["messages"].([]any); ok {
		input := convertMessagesToInput(messages)
		openaiReq["input"] = input
	}

	// 转换 stream
	if stream, ok := claudeReq["stream"].(bool); ok {
		openaiReq["stream"] = stream
	} else {
		openaiReq["stream"] = true // OpenAI 默认使用流式
	}

	// 设置 store = false (OpenAI 要求)
	openaiReq["store"] = false

	return json.Marshal(openaiReq)
}

// convertMessagesToInput 将 Claude messages 数组转换为 OpenAI input 数组
func convertMessagesToInput(messages []any) []any {
	input := make([]any, 0, len(messages))

	for _, msg := range messages {
		msgMap, ok := msg.(map[string]any)
		if !ok {
			continue
		}

		role, _ := msgMap["role"].(string)
		content := msgMap["content"]

		// 处理不同类型的 content
		switch c := content.(type) {
		case string:
			// 简单文本消息
			input = append(input, map[string]any{
				"type": "message",
				"role": role,
				"content": []any{
					map[string]any{
						"type": "input_text",
						"text": c,
					},
				},
			})
		case []any:
			// 复杂内容（可能包含文本、图片等）
			convertedContent := convertContentBlocks(c)
			input = append(input, map[string]any{
				"type":    "message",
				"role":    role,
				"content": convertedContent,
			})
		}
	}

	return input
}

// convertContentBlocks 转换 Claude content blocks 到 OpenAI 格式
func convertContentBlocks(blocks []any) []any {
	result := make([]any, 0, len(blocks))

	for _, block := range blocks {
		blockMap, ok := block.(map[string]any)
		if !ok {
			continue
		}

		blockType, _ := blockMap["type"].(string)

		switch blockType {
		case "text":
			if text, ok := blockMap["text"].(string); ok {
				result = append(result, map[string]any{
					"type": "input_text",
					"text": text,
				})
			}
		case "image":
			// 转换图片格式
			if source, ok := blockMap["source"].(map[string]any); ok {
				if data, ok := source["data"].(string); ok {
					mediaType, _ := source["media_type"].(string)
					result = append(result, map[string]any{
						"type": "input_image",
						"image_url": map[string]any{
							"url": "data:" + mediaType + ";base64," + data,
						},
					})
				}
			}
		default:
			// 其他类型直接保留
			result = append(result, blockMap)
		}
	}

	return result
}

// OpenAIToClaudeResponse 将 OpenAI SSE 响应转换为 Claude SSE 响应格式
// 这个函数处理单个 SSE 事件
func OpenAIToClaudeResponse(openaiEvent string, requestModel string) string {
	// 移除 "data: " 前缀
	data := strings.TrimPrefix(openaiEvent, "data: ")
	data = strings.TrimSpace(data)

	if data == "" || data == "[DONE]" {
		return ""
	}

	var event map[string]any
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return ""
	}

	eventType, _ := event["type"].(string)

	switch eventType {
	case "response.output_text.delta":
		// 文本增量
		if delta, ok := event["delta"].(string); ok {
			claudeEvent := map[string]any{
				"type": "content_block_delta",
				"index": 0,
				"delta": map[string]any{
					"type": "text_delta",
					"text": delta,
				},
			}
			jsonBytes, _ := json.Marshal(claudeEvent)
			return "event: content_block_delta\ndata: " + string(jsonBytes) + "\n\n"
		}

	case "response.completed":
		// 响应完成
		usage := extractOpenAIUsage(event)

		// 发送 message_stop 事件
		stopEvent := map[string]any{
			"type": "message_stop",
		}
		stopBytes, _ := json.Marshal(stopEvent)

		// 发送最终的 message_delta 包含 usage
		deltaEvent := map[string]any{
			"type": "message_delta",
			"delta": map[string]any{
				"stop_reason": "end_turn",
			},
			"usage": usage,
		}
		deltaBytes, _ := json.Marshal(deltaEvent)

		return "event: message_delta\ndata: " + string(deltaBytes) + "\n\n" +
			"event: message_stop\ndata: " + string(stopBytes) + "\n\n"

	case "response.output_text.done":
		// 文本块完成
		claudeEvent := map[string]any{
			"type":  "content_block_stop",
			"index": 0,
		}
		jsonBytes, _ := json.Marshal(claudeEvent)
		return "event: content_block_stop\ndata: " + string(jsonBytes) + "\n\n"
	}

	return ""
}

// extractOpenAIUsage 从 OpenAI 响应中提取 usage 信息
func extractOpenAIUsage(event map[string]any) map[string]any {
	usage := map[string]any{
		"input_tokens":  0,
		"output_tokens": 0,
	}

	if response, ok := event["response"].(map[string]any); ok {
		if u, ok := response["usage"].(map[string]any); ok {
			if inputTokens, ok := u["input_tokens"].(float64); ok {
				usage["input_tokens"] = int(inputTokens)
			}
			if outputTokens, ok := u["output_tokens"].(float64); ok {
				usage["output_tokens"] = int(outputTokens)
			}
		}
	}

	return usage
}

// GenerateClaudeMessageStart 生成 Claude 格式的 message_start 事件
func GenerateClaudeMessageStart(model string) string {
	event := map[string]any{
		"type": "message_start",
		"message": map[string]any{
			"id":           "msg_openai_compat",
			"type":         "message",
			"role":         "assistant",
			"content":      []any{},
			"model":        model,
			"stop_reason":  nil,
			"stop_sequence": nil,
			"usage": map[string]any{
				"input_tokens":  0,
				"output_tokens": 0,
			},
		},
	}
	jsonBytes, _ := json.Marshal(event)
	return "event: message_start\ndata: " + string(jsonBytes) + "\n\n"
}

// GenerateClaudeContentBlockStart 生成 Claude 格式的 content_block_start 事件
func GenerateClaudeContentBlockStart() string {
	event := map[string]any{
		"type":  "content_block_start",
		"index": 0,
		"content_block": map[string]any{
			"type": "text",
			"text": "",
		},
	}
	jsonBytes, _ := json.Marshal(event)
	return "event: content_block_start\ndata: " + string(jsonBytes) + "\n\n"
}
