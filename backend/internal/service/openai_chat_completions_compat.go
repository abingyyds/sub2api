package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/gin-gonic/gin"
)

const openAIChatCompletionsModeContextKey = "openai_chat_completions_mode"

// EnableOpenAIChatCompletionsMode marks the current request as chat-completions compatible.
func EnableOpenAIChatCompletionsMode(c *gin.Context) {
	if c != nil {
		c.Set(openAIChatCompletionsModeContextKey, true)
	}
}

func isOpenAIChatCompletionsMode(c *gin.Context) bool {
	if c == nil {
		return false
	}
	value, ok := c.Get(openAIChatCompletionsModeContextKey)
	flag, _ := value.(bool)
	return ok && flag
}

// ConvertChatCompletionsRequest normalizes a Chat Completions request into a Responses request.
func ConvertChatCompletionsRequest(reqBody map[string]any) (bool, error) {
	if reqBody == nil {
		return false, fmt.Errorf("request body is nil")
	}

	modified := false

	if rawMaxTokens, ok := reqBody["max_tokens"]; ok {
		if _, hasMaxOutputTokens := reqBody["max_output_tokens"]; !hasMaxOutputTokens {
			reqBody["max_output_tokens"] = rawMaxTokens
		}
		delete(reqBody, "max_tokens")
		modified = true
	}
	if rawMaxCompletionTokens, ok := reqBody["max_completion_tokens"]; ok {
		if _, hasMaxOutputTokens := reqBody["max_output_tokens"]; !hasMaxOutputTokens {
			reqBody["max_output_tokens"] = rawMaxCompletionTokens
		}
		delete(reqBody, "max_completion_tokens")
		modified = true
	}

	if normalizeCodexTools(reqBody) {
		modified = true
	}
	if normalizeChatCompletionsToolChoice(reqBody) {
		modified = true
	}

	rawMessages, hasMessages := reqBody["messages"]
	if !hasMessages || rawMessages == nil {
		return modified, nil
	}

	messages, ok := rawMessages.([]any)
	if !ok {
		return false, fmt.Errorf("messages must be an array")
	}

	input, instructions, err := convertChatMessagesToResponsesInput(messages)
	if err != nil {
		return false, err
	}

	if _, hasInput := reqBody["input"]; !hasInput || reqBody["input"] == nil {
		reqBody["input"] = input
		modified = true
	}

	delete(reqBody, "messages")
	modified = true

	if instructions != "" {
		existingInstructions, _ := reqBody["instructions"].(string)
		existingInstructions = strings.TrimSpace(existingInstructions)
		switch {
		case existingInstructions == "":
			reqBody["instructions"] = instructions
		case existingInstructions != instructions:
			reqBody["instructions"] = existingInstructions + "\n\n" + instructions
		}
		modified = true
	}

	return modified, nil
}

func normalizeChatCompletionsToolChoice(reqBody map[string]any) bool {
	raw, exists := reqBody["tool_choice"]
	if !exists || raw == nil {
		return false
	}

	toolChoice, ok := raw.(map[string]any)
	if !ok {
		return false
	}

	functionValue, hasFunction := toolChoice["function"]
	function, ok := functionValue.(map[string]any)
	if !hasFunction || functionValue == nil || !ok {
		return false
	}

	name, _ := function["name"].(string)
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}

	reqBody["tool_choice"] = map[string]any{
		"type": "function",
		"name": name,
	}
	return true
}

func convertChatMessagesToResponsesInput(messages []any) ([]any, string, error) {
	input := make([]any, 0, len(messages))
	systemInstructions := make([]string, 0, 2)

	for _, rawMessage := range messages {
		message, ok := rawMessage.(map[string]any)
		if !ok {
			return nil, "", fmt.Errorf("message must be an object")
		}

		role, _ := message["role"].(string)
		role = strings.TrimSpace(role)
		if role == "" {
			return nil, "", fmt.Errorf("message role is required")
		}

		switch role {
		case "system", "developer":
			text := flattenChatMessageContent(message["content"])
			if text != "" {
				systemInstructions = append(systemInstructions, text)
			}
			continue

		case "tool":
			callID, _ := message["tool_call_id"].(string)
			if strings.TrimSpace(callID) == "" {
				callID, _ = message["call_id"].(string)
			}
			callID = strings.TrimSpace(callID)
			if callID == "" {
				return nil, "", fmt.Errorf("tool message requires tool_call_id")
			}
			input = append(input, map[string]any{
				"type":    "function_call_output",
				"call_id": callID,
				"output":  flattenChatMessageContent(message["content"]),
			})
			continue
		}

		convertedContent, err := convertChatMessageContentToResponsesContent(message["content"])
		if err != nil {
			return nil, "", err
		}

		messageItem := map[string]any{
			"role": role,
		}
		if convertedContent != nil {
			messageItem["content"] = convertedContent
		}
		input = append(input, messageItem)

		if role == "assistant" {
			if toolCalls, ok := message["tool_calls"].([]any); ok {
				convertedToolCalls, err := convertAssistantToolCallsToResponses(toolCalls)
				if err != nil {
					return nil, "", err
				}
				input = append(input, convertedToolCalls...)
			}
		}
	}

	return input, strings.Join(systemInstructions, "\n\n"), nil
}

func convertAssistantToolCallsToResponses(toolCalls []any) ([]any, error) {
	converted := make([]any, 0, len(toolCalls))
	for _, rawToolCall := range toolCalls {
		toolCall, ok := rawToolCall.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("tool_calls item must be an object")
		}

		callID, _ := toolCall["id"].(string)
		callID = strings.TrimSpace(callID)

		functionValue, ok := toolCall["function"].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("tool_calls.function must be an object")
		}

		name, _ := functionValue["name"].(string)
		name = strings.TrimSpace(name)
		if name == "" {
			return nil, fmt.Errorf("tool call function name is required")
		}

		arguments, _ := functionValue["arguments"].(string)
		item := map[string]any{
			"type":      "function_call",
			"name":      name,
			"arguments": arguments,
		}
		if callID != "" {
			item["call_id"] = callID
			item["id"] = callID
		}
		converted = append(converted, item)
	}

	return converted, nil
}

func convertChatMessageContentToResponsesContent(content any) (any, error) {
	switch value := content.(type) {
	case nil:
		return nil, nil
	case string:
		return value, nil
	case []any:
		parts := make([]any, 0, len(value))
		for _, rawPart := range value {
			part, ok := rawPart.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("message content part must be an object")
			}

			partType, _ := part["type"].(string)
			partType = strings.TrimSpace(partType)

			switch partType {
			case "text", "input_text":
				text, _ := part["text"].(string)
				parts = append(parts, map[string]any{
					"type": "input_text",
					"text": text,
				})
			case "image_url":
				imageURL, ok := part["image_url"].(map[string]any)
				if !ok {
					return nil, fmt.Errorf("image_url part must include image_url.url")
				}
				urlValue, _ := imageURL["url"].(string)
				urlValue = strings.TrimSpace(urlValue)
				if urlValue == "" {
					return nil, fmt.Errorf("image_url.url is required")
				}
				converted := map[string]any{
					"type":      "input_image",
					"image_url": urlValue,
				}
				if detail, ok := imageURL["detail"].(string); ok && strings.TrimSpace(detail) != "" {
					converted["detail"] = detail
				}
				parts = append(parts, converted)
			case "input_image":
				parts = append(parts, part)
			default:
				if text, ok := part["text"].(string); ok {
					parts = append(parts, map[string]any{
						"type": "input_text",
						"text": text,
					})
					continue
				}
				return nil, fmt.Errorf("unsupported message content part type: %s", partType)
			}
		}
		return parts, nil
	default:
		return flattenChatMessageContent(content), nil
	}
}

func flattenChatMessageContent(content any) string {
	switch value := content.(type) {
	case nil:
		return ""
	case string:
		return value
	case []any:
		var parts []string
		for _, rawPart := range value {
			part, ok := rawPart.(map[string]any)
			if !ok {
				continue
			}
			if text, ok := part["text"].(string); ok && text != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "")
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return ""
		}
		return string(data)
	}
}

func (s *OpenAIGatewayService) convertResponsesBodyToChatCompletions(body []byte, fallbackModel string) ([]byte, *OpenAIUsage, error) {
	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}

	usage := extractOpenAIUsageFromResponsesObject(response)
	chat := buildChatCompletionsResponseFromResponsesObject(response, fallbackModel, usage)
	encoded, err := json.Marshal(chat)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal chat completion response: %w", err)
	}
	return encoded, usage, nil
}

func extractOpenAIUsageFromResponsesObject(response map[string]any) *OpenAIUsage {
	usage := &OpenAIUsage{}
	rawUsage, ok := response["usage"].(map[string]any)
	if !ok {
		return usage
	}

	usage.InputTokens = intValue(rawUsage["input_tokens"])
	usage.OutputTokens = intValue(rawUsage["output_tokens"])
	if details, ok := rawUsage["input_tokens_details"].(map[string]any); ok {
		usage.CacheReadInputTokens = intValue(details["cached_tokens"])
	}
	return usage
}

func buildChatCompletionsResponseFromResponsesObject(response map[string]any, fallbackModel string, usage *OpenAIUsage) map[string]any {
	model, _ := response["model"].(string)
	if model == "" {
		model = fallbackModel
	}

	id, _ := response["id"].(string)
	if id == "" {
		id = "chatcmpl-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	created := parseResponseCreatedAt(response)
	message, finishReason := extractChatMessageFromResponsesOutput(response["output"])

	totalTokens := usage.InputTokens + usage.OutputTokens
	result := map[string]any{
		"id":      id,
		"object":  "chat.completion",
		"created": created,
		"model":   model,
		"choices": []any{
			map[string]any{
				"index":         0,
				"message":       message,
				"finish_reason": finishReason,
			},
		},
		"usage": map[string]any{
			"prompt_tokens":     usage.InputTokens,
			"completion_tokens": usage.OutputTokens,
			"total_tokens":      totalTokens,
		},
	}
	return result
}

func extractChatMessageFromResponsesOutput(rawOutput any) (map[string]any, string) {
	message := map[string]any{
		"role": "assistant",
	}
	var contentParts []string
	var toolCalls []any

	output, ok := rawOutput.([]any)
	if !ok {
		message["content"] = ""
		return message, "stop"
	}

	for _, rawItem := range output {
		item, ok := rawItem.(map[string]any)
		if !ok {
			continue
		}

		itemType, _ := item["type"].(string)
		switch itemType {
		case "message":
			if role, ok := item["role"].(string); ok && strings.TrimSpace(role) != "" {
				message["role"] = role
			}
			contentParts = append(contentParts, extractTextFromResponsesMessageContent(item["content"]))
		case "function_call":
			toolCall := map[string]any{
				"type": "function",
				"function": map[string]any{
					"name":      stringValue(item["name"]),
					"arguments": stringValue(item["arguments"]),
				},
			}
			callID := stringValue(item["call_id"])
			if callID == "" {
				callID = stringValue(item["id"])
			}
			if callID != "" {
				toolCall["id"] = callID
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}

	content := strings.Join(contentParts, "")
	if len(toolCalls) > 0 {
		message["tool_calls"] = toolCalls
		if strings.TrimSpace(content) == "" {
			message["content"] = nil
		} else {
			message["content"] = content
		}
		return message, "tool_calls"
	}

	message["content"] = content
	return message, "stop"
}

func extractTextFromResponsesMessageContent(content any) string {
	switch value := content.(type) {
	case string:
		return value
	case []any:
		var builder strings.Builder
		for _, rawPart := range value {
			part, ok := rawPart.(map[string]any)
			if !ok {
				continue
			}
			partType, _ := part["type"].(string)
			switch partType {
			case "output_text", "input_text", "text":
				builder.WriteString(stringValue(part["text"]))
			}
		}
		return builder.String()
	default:
		return ""
	}
}

func parseResponseCreatedAt(response map[string]any) int64 {
	if createdAt, ok := response["created_at"].(string); ok {
		if ts, err := time.Parse(time.RFC3339, createdAt); err == nil {
			return ts.Unix()
		}
	}
	return time.Now().Unix()
}

type chatCompletionsStreamState struct {
	ID            string
	Model         string
	Created       int64
	SentRole      bool
	Finished      bool
	HasToolCalls  bool
	ToolCallIndex map[string]int
}

func newChatCompletionsStreamState(fallbackID, fallbackModel string) *chatCompletionsStreamState {
	if fallbackID == "" {
		fallbackID = "chatcmpl-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return &chatCompletionsStreamState{
		ID:            fallbackID,
		Model:         fallbackModel,
		Created:       time.Now().Unix(),
		ToolCallIndex: make(map[string]int),
	}
}

func (s *OpenAIGatewayService) convertResponsesStreamEventToChatCompletions(data string, state *chatCompletionsStreamState) ([]string, error) {
	if data == "" {
		return nil, nil
	}
	if data == "[DONE]" {
		if state.Finished {
			return []string{"data: [DONE]\n\n"}, nil
		}
		return nil, nil
	}

	var event map[string]any
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return nil, nil
	}

	if response, ok := event["response"].(map[string]any); ok {
		if id := stringValue(response["id"]); id != "" {
			state.ID = id
		}
		if model := stringValue(response["model"]); model != "" {
			state.Model = model
		}
	}
	if item, ok := event["item"].(map[string]any); ok {
		if itemType := stringValue(item["type"]); itemType == "function_call" {
			state.HasToolCalls = true
		}
	}

	eventType := stringValue(event["type"])
	switch eventType {
	case "response.output_text.delta":
		delta := stringValue(event["delta"])
		if delta == "" {
			return nil, nil
		}
		lines := make([]string, 0, 2)
		if !state.SentRole {
			lines = append(lines, encodeChatCompletionsChunk(state, map[string]any{"role": "assistant"}, nil))
			state.SentRole = true
		}
		lines = append(lines, encodeChatCompletionsChunk(state, map[string]any{"content": delta}, nil))
		return lines, nil

	case "response.output_item.added":
		item, _ := event["item"].(map[string]any)
		if stringValue(item["type"]) != "function_call" {
			return nil, nil
		}

		state.HasToolCalls = true
		callID := stringValue(item["call_id"])
		if callID == "" {
			callID = stringValue(item["id"])
		}
		index := len(state.ToolCallIndex)
		if existing, ok := state.ToolCallIndex[callID]; ok {
			index = existing
		} else {
			state.ToolCallIndex[callID] = index
		}

		chunk := map[string]any{
			"tool_calls": []any{
				map[string]any{
					"index": index,
					"id":    callID,
					"type":  "function",
					"function": map[string]any{
						"name":      stringValue(item["name"]),
						"arguments": "",
					},
				},
			},
		}
		return []string{encodeChatCompletionsChunk(state, chunk, nil)}, nil

	case "response.function_call_arguments.delta":
		callID := stringValue(event["call_id"])
		if callID == "" {
			callID = stringValue(event["item_id"])
		}
		index, ok := state.ToolCallIndex[callID]
		if !ok {
			index = len(state.ToolCallIndex)
			state.ToolCallIndex[callID] = index
		}
		chunk := map[string]any{
			"tool_calls": []any{
				map[string]any{
					"index": index,
					"function": map[string]any{
						"arguments": stringValue(event["delta"]),
					},
				},
			},
		}
		return []string{encodeChatCompletionsChunk(state, chunk, nil)}, nil

	case "response.completed", "response.done":
		if state.Finished {
			return nil, nil
		}
		state.Finished = true
		reason := "stop"
		if state.HasToolCalls {
			reason = "tool_calls"
		}
		return []string{
			encodeChatCompletionsChunk(state, map[string]any{}, &reason),
			"data: [DONE]\n\n",
		}, nil
	}

	return nil, nil
}

func encodeChatCompletionsChunk(state *chatCompletionsStreamState, delta map[string]any, finishReason *string) string {
	chunk := map[string]any{
		"id":      state.ID,
		"object":  "chat.completion.chunk",
		"created": state.Created,
		"model":   state.Model,
		"choices": []any{
			map[string]any{
				"index":         0,
				"delta":         delta,
				"finish_reason": finishReason,
			},
		},
	}
	data, _ := json.Marshal(chunk)
	return "data: " + string(data) + "\n\n"
}

func (s *OpenAIGatewayService) handleChatCompletionsStreamingResponse(_ context.Context, resp *http.Response, c *gin.Context, _ *Account, startTime time.Time, originalModel, mappedModel string) (*openaiStreamingResult, error) {
	if s.cfg != nil {
		responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.cfg.Security.ResponseHeaders)
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	if v := resp.Header.Get("x-request-id"); v != "" {
		c.Header("x-request-id", v)
	}

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming not supported")
	}

	usage := &OpenAIUsage{}
	var firstTokenMs *int
	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.cfg != nil && s.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.cfg.Gateway.MaxLineSize
	}
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	state := newChatCompletionsStreamState(resp.Header.Get("x-request-id"), originalModel)
	if originalModel != "" {
		state.Model = originalModel
	}

	needModelReplace := originalModel != mappedModel

	for scanner.Scan() {
		line := scanner.Text()
		if !openaiSSEDataRe.MatchString(line) {
			if strings.HasPrefix(line, ":") {
				if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
					return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs}, err
				}
				flusher.Flush()
			}
			continue
		}

		data := openaiSSEDataRe.ReplaceAllString(line, "")
		if needModelReplace {
			line = s.replaceModelInSSELine(line, mappedModel, originalModel)
			data = openaiSSEDataRe.ReplaceAllString(line, "")
		}
		if correctedData, corrected := s.toolCorrector.CorrectToolCallsInSSEData(data); corrected {
			data = correctedData
		}

		if firstTokenMs == nil && data != "" && data != "[DONE]" {
			ms := int(time.Since(startTime).Milliseconds())
			firstTokenMs = &ms
		}
		s.parseSSEUsage(data, usage)

		lines, err := s.convertResponsesStreamEventToChatCompletions(data, state)
		if err != nil {
			return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs}, err
		}
		for _, out := range lines {
			if _, err := fmt.Fprint(w, out); err != nil {
				return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs}, err
			}
			flusher.Flush()
		}
	}

	if err := scanner.Err(); err != nil {
		return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs}, err
	}

	if !state.Finished {
		reason := "stop"
		if state.HasToolCalls {
			reason = "tool_calls"
		}
		if _, err := fmt.Fprint(w, encodeChatCompletionsChunk(state, map[string]any{}, &reason)); err != nil {
			return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs}, err
		}
		if _, err := fmt.Fprint(w, "data: [DONE]\n\n"); err != nil {
			return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs}, err
		}
		flusher.Flush()
	}

	return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs}, nil
}

func (s *OpenAIGatewayService) handleChatCompletionsNonStreamingResponse(_ context.Context, resp *http.Response, c *gin.Context, account *Account, originalModel, mappedModel string) (*OpenAIUsage, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if account.Type == AccountTypeOAuth {
		bodyLooksLikeSSE := bytes.Contains(body, []byte("data:")) || bytes.Contains(body, []byte("event:"))
		if isEventStreamResponse(resp.Header) || bodyLooksLikeSSE {
			bodyText := string(body)
			finalResponse, ok := extractCodexFinalResponse(bodyText)
			if !ok {
				return nil, fmt.Errorf("failed to extract final response from OAuth SSE body")
			}
			body = finalResponse
		}
	}

	if originalModel != mappedModel {
		body = s.replaceModelInResponseBody(body, mappedModel, originalModel)
	}
	body = s.correctToolCallsInResponseBody(body)

	chatBody, usage, err := s.convertResponsesBodyToChatCompletions(body, originalModel)
	if err != nil {
		return nil, err
	}

	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.cfg.Security.ResponseHeaders)
	c.Data(resp.StatusCode, "application/json", chatBody)
	return usage, nil
}

func intValue(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		i, _ := v.Int64()
		return int(i)
	default:
		return 0
	}
}

func stringValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(bytes.Trim(data, `"`))
	}
}
