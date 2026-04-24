package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const DefaultOpenAIImageModel = "gpt-image-2"

type OpenAIImagesUpload struct {
	FieldName   string
	FileName    string
	ContentType string
	Data        []byte
}

type OpenAIImagesRequest struct {
	Endpoint          string
	ContentType       string
	Multipart         bool
	Model             string
	Prompt            string
	Stream            bool
	N                 int
	Size              string
	ResponseFormat    string
	Quality           string
	Background        string
	OutputFormat      string
	Moderation        string
	Style             string
	InputImageURLs    []string
	MaskImageURL      string
	Uploads           []OpenAIImagesUpload
	MaskUpload        *OpenAIImagesUpload
	OutputCompression *int
	PartialImages     *int
}

type openAIResponsesImageResult struct {
	Result        string
	RevisedPrompt string
	OutputFormat  string
	Size          string
	Background    string
	Quality       string
	Model         string
}

func (r *OpenAIImagesRequest) IsEdits() bool {
	return r != nil && r.Endpoint == OpenAIEndpointImagesEdits
}

func ParseOpenAIImagesRequest(body []byte, contentType, endpoint string) (*OpenAIImagesRequest, error) {
	req := &OpenAIImagesRequest{
		Endpoint:    endpoint,
		ContentType: strings.TrimSpace(contentType),
		N:           1,
	}

	mediaType := strings.ToLower(strings.TrimSpace(contentType))
	if mediaType != "" {
		parsedMediaType, _, err := mime.ParseMediaType(contentType)
		if err == nil {
			mediaType = strings.ToLower(parsedMediaType)
		}
	}

	if strings.HasPrefix(mediaType, "multipart/form-data") {
		req.Multipart = true
		if err := parseOpenAIImagesMultipartRequest(body, contentType, req); err != nil {
			return nil, err
		}
	} else {
		if len(body) == 0 || !gjson.ValidBytes(body) {
			return nil, fmt.Errorf("failed to parse request body")
		}
		if err := parseOpenAIImagesJSONRequest(body, req); err != nil {
			return nil, err
		}
	}

	if strings.TrimSpace(req.Model) == "" {
		req.Model = DefaultOpenAIImageModel
	}
	if req.N <= 0 {
		req.N = 1
	}
	return req, validateOpenAIImagesModel(req.Model)
}

func parseOpenAIImagesJSONRequest(body []byte, req *OpenAIImagesRequest) error {
	req.Model = strings.TrimSpace(gjson.GetBytes(body, "model").String())
	req.Prompt = strings.TrimSpace(gjson.GetBytes(body, "prompt").String())
	req.Stream = gjson.GetBytes(body, "stream").Bool()

	if nResult := gjson.GetBytes(body, "n"); nResult.Exists() {
		if nResult.Type != gjson.Number {
			return fmt.Errorf("invalid n field type")
		}
		req.N = int(nResult.Int())
	}

	req.Size = strings.TrimSpace(gjson.GetBytes(body, "size").String())
	req.ResponseFormat = strings.ToLower(strings.TrimSpace(gjson.GetBytes(body, "response_format").String()))
	req.Quality = strings.TrimSpace(gjson.GetBytes(body, "quality").String())
	req.Background = strings.TrimSpace(gjson.GetBytes(body, "background").String())
	req.OutputFormat = strings.TrimSpace(gjson.GetBytes(body, "output_format").String())
	req.Moderation = strings.TrimSpace(gjson.GetBytes(body, "moderation").String())
	req.Style = strings.TrimSpace(gjson.GetBytes(body, "style").String())

	if outputCompression := gjson.GetBytes(body, "output_compression"); outputCompression.Exists() {
		if outputCompression.Type != gjson.Number {
			return fmt.Errorf("invalid output_compression field type")
		}
		value := int(outputCompression.Int())
		req.OutputCompression = &value
	}
	if partialImages := gjson.GetBytes(body, "partial_images"); partialImages.Exists() {
		if partialImages.Type != gjson.Number {
			return fmt.Errorf("invalid partial_images field type")
		}
		value := int(partialImages.Int())
		req.PartialImages = &value
	}

	if req.IsEdits() {
		if images := gjson.GetBytes(body, "images"); images.Exists() && images.IsArray() {
			for _, item := range images.Array() {
				if imageURL := strings.TrimSpace(item.Get("image_url").String()); imageURL != "" {
					req.InputImageURLs = append(req.InputImageURLs, imageURL)
				}
			}
		}
		req.MaskImageURL = strings.TrimSpace(gjson.GetBytes(body, "mask.image_url").String())
	}

	return nil
}

func parseOpenAIImagesMultipartRequest(body []byte, contentType string, req *OpenAIImagesRequest) error {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("invalid multipart content-type: %w", err)
	}
	boundary := strings.TrimSpace(params["boundary"])
	if boundary == "" {
		return fmt.Errorf("multipart boundary is required")
	}

	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read multipart request: %w", err)
		}

		data, readErr := io.ReadAll(part)
		_ = part.Close()
		if readErr != nil {
			return fmt.Errorf("read multipart part: %w", readErr)
		}

		if part.FileName() == "" {
			switch strings.TrimSpace(part.FormName()) {
			case "model":
				req.Model = strings.TrimSpace(string(data))
			case "prompt":
				req.Prompt = strings.TrimSpace(string(data))
			case "stream":
				req.Stream = parseMultipartBool(string(data))
			case "n":
				if parsed := gjson.ParseBytes(data); parsed.Exists() && parsed.Type == gjson.Number {
					req.N = int(parsed.Int())
				}
			case "size":
				req.Size = strings.TrimSpace(string(data))
			case "response_format":
				req.ResponseFormat = strings.ToLower(strings.TrimSpace(string(data)))
			case "quality":
				req.Quality = strings.TrimSpace(string(data))
			case "background":
				req.Background = strings.TrimSpace(string(data))
			case "output_format":
				req.OutputFormat = strings.TrimSpace(string(data))
			case "moderation":
				req.Moderation = strings.TrimSpace(string(data))
			case "style":
				req.Style = strings.TrimSpace(string(data))
			case "output_compression":
				if parsed := gjson.ParseBytes(data); parsed.Exists() && parsed.Type == gjson.Number {
					value := int(parsed.Int())
					req.OutputCompression = &value
				}
			case "partial_images":
				if parsed := gjson.ParseBytes(data); parsed.Exists() && parsed.Type == gjson.Number {
					value := int(parsed.Int())
					req.PartialImages = &value
				}
			}
			continue
		}

		upload := OpenAIImagesUpload{
			FieldName:   part.FormName(),
			FileName:    part.FileName(),
			ContentType: strings.TrimSpace(part.Header.Get("Content-Type")),
			Data:        data,
		}
		if upload.FieldName == "mask" {
			req.MaskUpload = &upload
			continue
		}
		req.Uploads = append(req.Uploads, upload)
	}

	return nil
}

func isOpenAIImageGenerationModel(model string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(model)), "gpt-image-")
}

func validateOpenAIImagesModel(model string) error {
	if isOpenAIImageGenerationModel(model) {
		return nil
	}
	return fmt.Errorf("images endpoint requires an image model, got %q", model)
}

func openAIImageUploadToDataURL(upload OpenAIImagesUpload) (string, error) {
	if len(upload.Data) == 0 {
		return "", fmt.Errorf("upload %q is empty", strings.TrimSpace(upload.FileName))
	}
	contentType := strings.TrimSpace(upload.ContentType)
	if contentType == "" {
		contentType = http.DetectContentType(upload.Data)
	}
	return "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(upload.Data), nil
}

func buildOpenAIImagesResponsesRequest(parsed *OpenAIImagesRequest, toolModel string) ([]byte, error) {
	if parsed == nil {
		return nil, fmt.Errorf("parsed images request is required")
	}
	prompt := strings.TrimSpace(parsed.Prompt)
	if prompt == "" {
		return nil, fmt.Errorf("prompt is required")
	}

	inputImages := make([]string, 0, len(parsed.InputImageURLs)+len(parsed.Uploads))
	for _, imageURL := range parsed.InputImageURLs {
		if trimmed := strings.TrimSpace(imageURL); trimmed != "" {
			inputImages = append(inputImages, trimmed)
		}
	}
	for _, upload := range parsed.Uploads {
		dataURL, err := openAIImageUploadToDataURL(upload)
		if err != nil {
			return nil, err
		}
		inputImages = append(inputImages, dataURL)
	}
	if parsed.IsEdits() && len(inputImages) == 0 {
		return nil, fmt.Errorf("image input is required")
	}

	req := []byte(`{"instructions":"","stream":true,"reasoning":{"effort":"medium"},"parallel_tool_calls":true,"include":["reasoning.encrypted_content"],"model":"gpt-5.4-mini","store":false,"tool_choice":{"type":"image_generation"}}`)

	input := []byte(`[{"type":"message","role":"user","content":[{"type":"input_text","text":""}]}]`)
	input, _ = sjson.SetBytes(input, "0.content.0.text", prompt)
	for index, imageURL := range inputImages {
		part := []byte(`{"type":"input_image","image_url":""}`)
		part, _ = sjson.SetBytes(part, "image_url", imageURL)
		input, _ = sjson.SetRawBytes(input, fmt.Sprintf("0.content.%d", index+1), part)
	}
	req, _ = sjson.SetRawBytes(req, "input", input)

	action := "generate"
	if parsed.IsEdits() {
		action = "edit"
	}
	tool := []byte(`{"type":"image_generation","action":"","model":""}`)
	tool, _ = sjson.SetBytes(tool, "action", action)
	tool, _ = sjson.SetBytes(tool, "model", strings.TrimSpace(toolModel))

	for _, field := range []struct {
		path  string
		value string
	}{
		{path: "size", value: parsed.Size},
		{path: "quality", value: parsed.Quality},
		{path: "background", value: parsed.Background},
		{path: "output_format", value: parsed.OutputFormat},
		{path: "moderation", value: parsed.Moderation},
		{path: "style", value: parsed.Style},
	} {
		if trimmed := strings.TrimSpace(field.value); trimmed != "" {
			tool, _ = sjson.SetBytes(tool, field.path, trimmed)
		}
	}
	if parsed.OutputCompression != nil {
		tool, _ = sjson.SetBytes(tool, "output_compression", *parsed.OutputCompression)
	}
	if parsed.PartialImages != nil {
		tool, _ = sjson.SetBytes(tool, "partial_images", *parsed.PartialImages)
	}

	maskImageURL := strings.TrimSpace(parsed.MaskImageURL)
	if parsed.MaskUpload != nil {
		dataURL, err := openAIImageUploadToDataURL(*parsed.MaskUpload)
		if err != nil {
			return nil, err
		}
		maskImageURL = dataURL
	}
	if maskImageURL != "" {
		tool, _ = sjson.SetBytes(tool, "input_image_mask.image_url", maskImageURL)
	}

	req, _ = sjson.SetRawBytes(req, "tools", []byte(`[]`))
	req, _ = sjson.SetRawBytes(req, "tools.-1", tool)
	return req, nil
}

func extractOpenAISSEDataLine(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || !openaiSSEDataRe.MatchString(line) {
		return "", false
	}
	return openaiSSEDataRe.ReplaceAllString(line, ""), true
}

func extractOpenAIResponsesImageMetaFromLifecycleEvent(payload []byte) (openAIResponsesImageResult, int64, bool) {
	switch gjson.GetBytes(payload, "type").String() {
	case "response.created", "response.in_progress", "response.completed":
	default:
		return openAIResponsesImageResult{}, 0, false
	}

	response := gjson.GetBytes(payload, "response")
	if !response.Exists() {
		return openAIResponsesImageResult{}, 0, false
	}

	return openAIResponsesImageResult{
		OutputFormat: strings.TrimSpace(response.Get("tools.0.output_format").String()),
		Size:         strings.TrimSpace(response.Get("tools.0.size").String()),
		Background:   strings.TrimSpace(response.Get("tools.0.background").String()),
		Quality:      strings.TrimSpace(response.Get("tools.0.quality").String()),
		Model:        strings.TrimSpace(response.Get("tools.0.model").String()),
	}, response.Get("created_at").Int(), true
}

func mergeOpenAIResponsesImageMeta(dst *openAIResponsesImageResult, src openAIResponsesImageResult) {
	if dst == nil {
		return
	}
	if trimmed := strings.TrimSpace(src.OutputFormat); trimmed != "" {
		dst.OutputFormat = trimmed
	}
	if trimmed := strings.TrimSpace(src.Size); trimmed != "" {
		dst.Size = trimmed
	}
	if trimmed := strings.TrimSpace(src.Background); trimmed != "" {
		dst.Background = trimmed
	}
	if trimmed := strings.TrimSpace(src.Quality); trimmed != "" {
		dst.Quality = trimmed
	}
	if trimmed := strings.TrimSpace(src.Model); trimmed != "" {
		dst.Model = trimmed
	}
}

func extractOpenAIImagesFromResponsesCompleted(payload []byte) ([]openAIResponsesImageResult, int64, []byte, openAIResponsesImageResult, error) {
	if gjson.GetBytes(payload, "type").String() != "response.completed" {
		return nil, 0, nil, openAIResponsesImageResult{}, fmt.Errorf("unexpected event type")
	}

	createdAt := gjson.GetBytes(payload, "response.created_at").Int()
	if createdAt <= 0 {
		createdAt = time.Now().Unix()
	}

	var (
		results   []openAIResponsesImageResult
		firstMeta openAIResponsesImageResult
	)
	output := gjson.GetBytes(payload, "response.output")
	if output.IsArray() {
		for _, item := range output.Array() {
			if item.Get("type").String() != "image_generation_call" {
				continue
			}
			result := strings.TrimSpace(item.Get("result").String())
			if result == "" {
				continue
			}
			entry := openAIResponsesImageResult{
				Result:        result,
				RevisedPrompt: strings.TrimSpace(item.Get("revised_prompt").String()),
				OutputFormat:  strings.TrimSpace(item.Get("output_format").String()),
				Size:          strings.TrimSpace(item.Get("size").String()),
				Background:    strings.TrimSpace(item.Get("background").String()),
				Quality:       strings.TrimSpace(item.Get("quality").String()),
				Model:         strings.TrimSpace(item.Get("model").String()),
			}
			if len(results) == 0 {
				firstMeta = entry
			}
			results = append(results, entry)
		}
	}

	var usageRaw []byte
	if usage := gjson.GetBytes(payload, "response.tool_usage.image_gen"); usage.Exists() && usage.IsObject() {
		usageRaw = []byte(usage.Raw)
	}
	return results, createdAt, usageRaw, firstMeta, nil
}

func extractOpenAIImageFromResponsesOutputItemDone(payload []byte) (openAIResponsesImageResult, string, bool, error) {
	if gjson.GetBytes(payload, "type").String() != "response.output_item.done" {
		return openAIResponsesImageResult{}, "", false, fmt.Errorf("unexpected event type")
	}

	item := gjson.GetBytes(payload, "item")
	if !item.Exists() || item.Get("type").String() != "image_generation_call" {
		return openAIResponsesImageResult{}, "", false, nil
	}

	result := strings.TrimSpace(item.Get("result").String())
	if result == "" {
		return openAIResponsesImageResult{}, "", false, nil
	}

	entry := openAIResponsesImageResult{
		Result:        result,
		RevisedPrompt: strings.TrimSpace(item.Get("revised_prompt").String()),
		OutputFormat:  strings.TrimSpace(item.Get("output_format").String()),
		Size:          strings.TrimSpace(item.Get("size").String()),
		Background:    strings.TrimSpace(item.Get("background").String()),
		Quality:       strings.TrimSpace(item.Get("quality").String()),
		Model:         strings.TrimSpace(item.Get("model").String()),
	}
	return entry, strings.TrimSpace(item.Get("id").String()), true, nil
}

func openAIResponsesImageResultKey(itemID string, result openAIResponsesImageResult) string {
	if trimmed := strings.TrimSpace(result.Result); trimmed != "" {
		return strings.TrimSpace(result.OutputFormat) + "|" + trimmed
	}
	return "item:" + strings.TrimSpace(itemID)
}

func appendOpenAIResponsesImageResultDedup(results *[]openAIResponsesImageResult, seen map[string]struct{}, itemID string, result openAIResponsesImageResult) bool {
	if results == nil {
		return false
	}
	key := openAIResponsesImageResultKey(itemID, result)
	if key != "" {
		if _, exists := seen[key]; exists {
			return false
		}
		seen[key] = struct{}{}
	}
	*results = append(*results, result)
	return true
}

func collectOpenAIImagesFromResponsesBody(body []byte) ([]openAIResponsesImageResult, int64, []byte, openAIResponsesImageResult, bool, error) {
	var (
		fallbackResults []openAIResponsesImageResult
		fallbackSeen    = make(map[string]struct{})
		createdAt       int64
		usageRaw        []byte
		foundFinal      bool
		responseMeta    openAIResponsesImageResult
	)

	for _, line := range bytes.Split(body, []byte("\n")) {
		line = bytes.TrimRight(line, "\r")
		data, ok := extractOpenAISSEDataLine(string(line))
		if !ok || data == "" || data == "[DONE]" {
			continue
		}
		payload := []byte(data)
		if !gjson.ValidBytes(payload) {
			continue
		}
		if meta, eventCreatedAt, ok := extractOpenAIResponsesImageMetaFromLifecycleEvent(payload); ok {
			mergeOpenAIResponsesImageMeta(&responseMeta, meta)
			if eventCreatedAt > 0 {
				createdAt = eventCreatedAt
			}
		}

		switch gjson.GetBytes(payload, "type").String() {
		case "response.output_item.done":
			result, itemID, ok, err := extractOpenAIImageFromResponsesOutputItemDone(payload)
			if err != nil {
				return nil, 0, nil, openAIResponsesImageResult{}, false, err
			}
			if ok {
				mergeOpenAIResponsesImageMeta(&result, responseMeta)
				appendOpenAIResponsesImageResultDedup(&fallbackResults, fallbackSeen, itemID, result)
			}
		case "response.completed":
			results, completedAt, completedUsageRaw, firstMeta, err := extractOpenAIImagesFromResponsesCompleted(payload)
			if err != nil {
				return nil, 0, nil, openAIResponsesImageResult{}, false, err
			}
			foundFinal = true
			if completedAt > 0 {
				createdAt = completedAt
			}
			if len(completedUsageRaw) > 0 {
				usageRaw = completedUsageRaw
			}
			if len(results) > 0 {
				mergeOpenAIResponsesImageMeta(&firstMeta, responseMeta)
				return results, createdAt, usageRaw, firstMeta, true, nil
			}
			if len(fallbackResults) > 0 {
				firstMeta = fallbackResults[0]
				mergeOpenAIResponsesImageMeta(&firstMeta, responseMeta)
				return fallbackResults, createdAt, usageRaw, firstMeta, true, nil
			}
		}
	}

	if len(fallbackResults) > 0 {
		firstMeta := fallbackResults[0]
		mergeOpenAIResponsesImageMeta(&firstMeta, responseMeta)
		return fallbackResults, createdAt, usageRaw, firstMeta, foundFinal, nil
	}
	return nil, createdAt, usageRaw, openAIResponsesImageResult{}, foundFinal, nil
}

func buildOpenAIImagesAPIResponse(results []openAIResponsesImageResult, createdAt int64, usageRaw []byte, firstMeta openAIResponsesImageResult, responseFormat string) ([]byte, error) {
	if createdAt <= 0 {
		createdAt = time.Now().Unix()
	}
	out := []byte(`{"created":0,"data":[]}`)
	out, _ = sjson.SetBytes(out, "created", createdAt)

	format := strings.ToLower(strings.TrimSpace(responseFormat))
	if format == "" {
		format = "b64_json"
	}
	for _, img := range results {
		item := []byte(`{}`)
		if format == "url" {
			item, _ = sjson.SetBytes(item, "url", "data:"+openAIImageOutputMIMEType(img.OutputFormat)+";base64,"+img.Result)
		} else {
			item, _ = sjson.SetBytes(item, "b64_json", img.Result)
		}
		if img.RevisedPrompt != "" {
			item, _ = sjson.SetBytes(item, "revised_prompt", img.RevisedPrompt)
		}
		out, _ = sjson.SetRawBytes(out, "data.-1", item)
	}
	if firstMeta.Background != "" {
		out, _ = sjson.SetBytes(out, "background", firstMeta.Background)
	}
	if firstMeta.OutputFormat != "" {
		out, _ = sjson.SetBytes(out, "output_format", firstMeta.OutputFormat)
	}
	if firstMeta.Quality != "" {
		out, _ = sjson.SetBytes(out, "quality", firstMeta.Quality)
	}
	if firstMeta.Size != "" {
		out, _ = sjson.SetBytes(out, "size", firstMeta.Size)
	}
	if firstMeta.Model != "" {
		out, _ = sjson.SetBytes(out, "model", firstMeta.Model)
	}
	if len(usageRaw) > 0 && gjson.ValidBytes(usageRaw) {
		out, _ = sjson.SetRawBytes(out, "usage", usageRaw)
	}
	return out, nil
}

func openAIImageOutputMIMEType(outputFormat string) string {
	if outputFormat == "" {
		return "image/png"
	}
	if strings.Contains(outputFormat, "/") {
		return outputFormat
	}
	switch strings.ToLower(strings.TrimSpace(outputFormat)) {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "webp":
		return "image/webp"
	default:
		return "image/png"
	}
}

func mergeOpenAIUsageRaw(usage *OpenAIUsage, usageRaw []byte) {
	if usage == nil || len(usageRaw) == 0 || !gjson.ValidBytes(usageRaw) {
		return
	}
	usage.InputTokens = int(gjson.GetBytes(usageRaw, "input_tokens").Int())
	usage.OutputTokens = int(gjson.GetBytes(usageRaw, "output_tokens").Int())
	usage.CacheReadInputTokens = int(gjson.GetBytes(usageRaw, "input_tokens_details.cached_tokens").Int())
}

func buildOpenAIImagesStreamPartialPayload(eventType, b64 string, partialImageIndex int64, responseFormat string, createdAt int64, meta openAIResponsesImageResult) []byte {
	if createdAt <= 0 {
		createdAt = time.Now().Unix()
	}
	payload := []byte(`{"type":"","created_at":0,"partial_image_index":0,"b64_json":""}`)
	payload, _ = sjson.SetBytes(payload, "type", eventType)
	payload, _ = sjson.SetBytes(payload, "created_at", createdAt)
	payload, _ = sjson.SetBytes(payload, "partial_image_index", partialImageIndex)
	payload, _ = sjson.SetBytes(payload, "b64_json", b64)
	if strings.EqualFold(strings.TrimSpace(responseFormat), "url") {
		payload, _ = sjson.SetBytes(payload, "url", "data:"+openAIImageOutputMIMEType(meta.OutputFormat)+";base64,"+b64)
	}
	if meta.Background != "" {
		payload, _ = sjson.SetBytes(payload, "background", meta.Background)
	}
	if meta.OutputFormat != "" {
		payload, _ = sjson.SetBytes(payload, "output_format", meta.OutputFormat)
	}
	if meta.Quality != "" {
		payload, _ = sjson.SetBytes(payload, "quality", meta.Quality)
	}
	if meta.Size != "" {
		payload, _ = sjson.SetBytes(payload, "size", meta.Size)
	}
	if meta.Model != "" {
		payload, _ = sjson.SetBytes(payload, "model", meta.Model)
	}
	return payload
}

func buildOpenAIImagesStreamCompletedPayload(eventType string, img openAIResponsesImageResult, responseFormat string, createdAt int64, usageRaw []byte) []byte {
	if createdAt <= 0 {
		createdAt = time.Now().Unix()
	}
	payload := []byte(`{"type":"","created_at":0,"b64_json":""}`)
	payload, _ = sjson.SetBytes(payload, "type", eventType)
	payload, _ = sjson.SetBytes(payload, "created_at", createdAt)
	payload, _ = sjson.SetBytes(payload, "b64_json", img.Result)
	if strings.EqualFold(strings.TrimSpace(responseFormat), "url") {
		payload, _ = sjson.SetBytes(payload, "url", "data:"+openAIImageOutputMIMEType(img.OutputFormat)+";base64,"+img.Result)
	}
	if img.Background != "" {
		payload, _ = sjson.SetBytes(payload, "background", img.Background)
	}
	if img.OutputFormat != "" {
		payload, _ = sjson.SetBytes(payload, "output_format", img.OutputFormat)
	}
	if img.Quality != "" {
		payload, _ = sjson.SetBytes(payload, "quality", img.Quality)
	}
	if img.Size != "" {
		payload, _ = sjson.SetBytes(payload, "size", img.Size)
	}
	if img.Model != "" {
		payload, _ = sjson.SetBytes(payload, "model", img.Model)
	}
	if len(usageRaw) > 0 && gjson.ValidBytes(usageRaw) {
		payload, _ = sjson.SetRawBytes(payload, "usage", usageRaw)
	}
	return payload
}

func openAIImagesStreamPrefix(parsed *OpenAIImagesRequest) string {
	if parsed != nil && parsed.IsEdits() {
		return "image_edit"
	}
	return "image_generation"
}

func buildOpenAIImagesStreamErrorBody(message string) []byte {
	payload := []byte(`{"type":"error","error":{"type":"upstream_error","message":""}}`)
	if strings.TrimSpace(message) == "" {
		message = "upstream request failed"
	}
	payload, _ = sjson.SetBytes(payload, "error.message", message)
	return payload
}

func (s *OpenAIGatewayService) writeOpenAIImagesStreamEvent(c *gin.Context, flusher http.Flusher, eventName string, payload []byte) error {
	if strings.TrimSpace(eventName) != "" {
		if _, err := fmt.Fprintf(c.Writer, "event: %s\n", eventName); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", payload); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

func (s *OpenAIGatewayService) handleOpenAIImagesOAuthNonStreamingResponse(resp *http.Response, c *gin.Context, responseFormat, fallbackModel string) (OpenAIUsage, int, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OpenAIUsage{}, 0, err
	}

	results, createdAt, usageRaw, firstMeta, _, err := collectOpenAIImagesFromResponsesBody(body)
	if err != nil {
		return OpenAIUsage{}, 0, err
	}
	if len(results) == 0 {
		return OpenAIUsage{}, 0, fmt.Errorf("upstream did not return image output")
	}
	if strings.TrimSpace(firstMeta.Model) == "" {
		firstMeta.Model = strings.TrimSpace(fallbackModel)
	}

	responseBody, err := buildOpenAIImagesAPIResponse(results, createdAt, usageRaw, firstMeta, responseFormat)
	if err != nil {
		return OpenAIUsage{}, 0, err
	}

	if s.cfg != nil {
		responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.cfg.Security.ResponseHeaders)
	}
	c.Data(resp.StatusCode, "application/json; charset=utf-8", responseBody)

	usage := OpenAIUsage{}
	mergeOpenAIUsageRaw(&usage, usageRaw)
	return usage, len(results), nil
}

func (s *OpenAIGatewayService) handleOpenAIImagesOAuthStreamingResponse(resp *http.Response, c *gin.Context, startTime time.Time, responseFormat, streamPrefix, fallbackModel string) (OpenAIUsage, int, *int, error) {
	if s.cfg != nil {
		responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.cfg.Security.ResponseHeaders)
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(resp.StatusCode)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return OpenAIUsage{}, 0, nil, fmt.Errorf("streaming is not supported by response writer")
	}

	format := strings.ToLower(strings.TrimSpace(responseFormat))
	if format == "" {
		format = "b64_json"
	}

	reader := bufio.NewReader(resp.Body)
	usage := OpenAIUsage{}
	imageCount := 0
	var firstTokenMs *int
	emitted := make(map[string]struct{})
	pendingResults := make([]openAIResponsesImageResult, 0, 1)
	pendingSeen := make(map[string]struct{})
	streamMeta := openAIResponsesImageResult{Model: strings.TrimSpace(fallbackModel)}
	var createdAt int64

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			trimmedLine := strings.TrimRight(string(line), "\r\n")
			data, ok := extractOpenAISSEDataLine(trimmedLine)
			if ok && data != "" && data != "[DONE]" {
				if firstTokenMs == nil {
					ms := int(time.Since(startTime).Milliseconds())
					firstTokenMs = &ms
				}
				dataBytes := []byte(data)
				if meta, eventCreatedAt, ok := extractOpenAIResponsesImageMetaFromLifecycleEvent(dataBytes); ok {
					mergeOpenAIResponsesImageMeta(&streamMeta, meta)
					if eventCreatedAt > 0 {
						createdAt = eventCreatedAt
					}
				}
				switch gjson.GetBytes(dataBytes, "type").String() {
				case "response.image_generation_call.partial_image":
					b64 := strings.TrimSpace(gjson.GetBytes(dataBytes, "partial_image_b64").String())
					if b64 != "" {
						partialMeta := streamMeta
						mergeOpenAIResponsesImageMeta(&partialMeta, openAIResponsesImageResult{
							OutputFormat: strings.TrimSpace(gjson.GetBytes(dataBytes, "output_format").String()),
							Background:   strings.TrimSpace(gjson.GetBytes(dataBytes, "background").String()),
						})
						payload := buildOpenAIImagesStreamPartialPayload(
							streamPrefix+".partial_image",
							b64,
							gjson.GetBytes(dataBytes, "partial_image_index").Int(),
							format,
							createdAt,
							partialMeta,
						)
						if writeErr := s.writeOpenAIImagesStreamEvent(c, flusher, streamPrefix+".partial_image", payload); writeErr != nil {
							return OpenAIUsage{}, imageCount, firstTokenMs, writeErr
						}
					}
				case "response.output_item.done":
					img, itemID, ok, extractErr := extractOpenAIImageFromResponsesOutputItemDone(dataBytes)
					if extractErr != nil {
						_ = s.writeOpenAIImagesStreamEvent(c, flusher, "error", buildOpenAIImagesStreamErrorBody(extractErr.Error()))
						return OpenAIUsage{}, imageCount, firstTokenMs, extractErr
					}
					if !ok {
						break
					}
					mergeOpenAIResponsesImageMeta(&streamMeta, img)
					mergeOpenAIResponsesImageMeta(&img, streamMeta)
					key := openAIResponsesImageResultKey(itemID, img)
					if _, exists := emitted[key]; exists {
						break
					}
					if _, exists := pendingSeen[key]; exists {
						break
					}
					pendingSeen[key] = struct{}{}
					pendingResults = append(pendingResults, img)
				case "response.completed":
					results, _, usageRaw, firstMeta, extractErr := extractOpenAIImagesFromResponsesCompleted(dataBytes)
					if extractErr != nil {
						_ = s.writeOpenAIImagesStreamEvent(c, flusher, "error", buildOpenAIImagesStreamErrorBody(extractErr.Error()))
						return OpenAIUsage{}, imageCount, firstTokenMs, extractErr
					}
					mergeOpenAIResponsesImageMeta(&streamMeta, firstMeta)
					finalResults := make([]openAIResponsesImageResult, 0, len(results)+len(pendingResults))
					finalSeen := make(map[string]struct{})
					for _, img := range results {
						mergeOpenAIResponsesImageMeta(&img, streamMeta)
						appendOpenAIResponsesImageResultDedup(&finalResults, finalSeen, "", img)
					}
					for _, img := range pendingResults {
						mergeOpenAIResponsesImageMeta(&img, streamMeta)
						appendOpenAIResponsesImageResultDedup(&finalResults, finalSeen, "", img)
					}
					if len(finalResults) == 0 {
						streamErr := fmt.Errorf("upstream did not return image output")
						_ = s.writeOpenAIImagesStreamEvent(c, flusher, "error", buildOpenAIImagesStreamErrorBody(streamErr.Error()))
						return OpenAIUsage{}, imageCount, firstTokenMs, streamErr
					}
					mergeOpenAIUsageRaw(&usage, usageRaw)
					for _, img := range finalResults {
						key := openAIResponsesImageResultKey("", img)
						if _, exists := emitted[key]; exists {
							continue
						}
						payload := buildOpenAIImagesStreamCompletedPayload(streamPrefix+".completed", img, format, createdAt, usageRaw)
						if writeErr := s.writeOpenAIImagesStreamEvent(c, flusher, streamPrefix+".completed", payload); writeErr != nil {
							return OpenAIUsage{}, imageCount, firstTokenMs, writeErr
						}
						emitted[key] = struct{}{}
					}
					imageCount = len(emitted)
					return usage, imageCount, firstTokenMs, nil
				case "error":
					streamErr := fmt.Errorf("%s", strings.TrimSpace(gjson.GetBytes(dataBytes, "error.message").String()))
					if streamErr.Error() == "" {
						streamErr = fmt.Errorf("upstream image generation failed")
					}
					_ = s.writeOpenAIImagesStreamEvent(c, flusher, "error", buildOpenAIImagesStreamErrorBody(streamErr.Error()))
					return OpenAIUsage{}, imageCount, firstTokenMs, streamErr
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			_ = s.writeOpenAIImagesStreamEvent(c, flusher, "error", buildOpenAIImagesStreamErrorBody(err.Error()))
			return OpenAIUsage{}, imageCount, firstTokenMs, err
		}
	}

	if len(pendingResults) > 0 {
		for _, img := range pendingResults {
			mergeOpenAIResponsesImageMeta(&img, streamMeta)
			key := openAIResponsesImageResultKey("", img)
			if _, exists := emitted[key]; exists {
				continue
			}
			payload := buildOpenAIImagesStreamCompletedPayload(streamPrefix+".completed", img, format, createdAt, nil)
			if writeErr := s.writeOpenAIImagesStreamEvent(c, flusher, streamPrefix+".completed", payload); writeErr != nil {
				return OpenAIUsage{}, imageCount, firstTokenMs, writeErr
			}
			emitted[key] = struct{}{}
		}
		imageCount = len(emitted)
		return usage, imageCount, firstTokenMs, nil
	}

	streamErr := fmt.Errorf("stream disconnected before image generation completed")
	_ = s.writeOpenAIImagesStreamEvent(c, flusher, "error", buildOpenAIImagesStreamErrorBody(streamErr.Error()))
	return OpenAIUsage{}, imageCount, firstTokenMs, streamErr
}

func (s *OpenAIGatewayService) forwardOpenAIImagesOAuth(ctx context.Context, c *gin.Context, account *Account, parsed *OpenAIImagesRequest, upstreamModel string) (*OpenAIForwardResult, error) {
	startTime := time.Now()
	requestModel := strings.TrimSpace(parsed.Model)
	if requestModel == "" {
		requestModel = DefaultOpenAIImageModel
	}
	if err := validateOpenAIImagesModel(requestModel); err != nil {
		return nil, err
	}

	if mapped := strings.TrimSpace(upstreamModel); mapped != "" {
		upstreamModel = mapped
	} else {
		upstreamModel = requestModel
	}

	responsesBody, err := buildOpenAIImagesResponsesRequest(parsed, upstreamModel)
	if err != nil {
		return nil, err
	}
	if c != nil {
		c.Set(OpsUpstreamRequestBodyKey, string(responsesBody))
	}

	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, err
	}
	upstreamReq, err := s.buildUpstreamRequest(ctx, c, account, responsesBody, token, true, parsed.Prompt, false, OpenAIEndpointResponses, "application/json")
	if err != nil {
		return nil, err
	}
	upstreamReq.Header.Set("Accept", "text/event-stream")

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}

	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		safeErr := sanitizeUpstreamErrorMessage(err.Error())
		setOpsUpstreamError(c, 0, safeErr, "")
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: 0,
			Kind:               "request_error",
			Message:            safeErr,
		})
		return nil, &UpstreamFailoverError{StatusCode: 0}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		if s.shouldFailoverUpstreamError(resp.StatusCode) {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewReader(respBody))
			upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
			upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  resp.Header.Get("x-request-id"),
				Kind:               "failover",
				Message:            upstreamMsg,
			})
			s.handleFailoverSideEffects(ctx, resp, account)
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCode}
		}
		return s.handleErrorResponse(ctx, resp, c, account)
	}

	var (
		usage        OpenAIUsage
		imageCount   int
		firstTokenMs *int
	)
	if parsed.Stream {
		usage, imageCount, firstTokenMs, err = s.handleOpenAIImagesOAuthStreamingResponse(resp, c, startTime, parsed.ResponseFormat, openAIImagesStreamPrefix(parsed), requestModel)
	} else {
		usage, imageCount, err = s.handleOpenAIImagesOAuthNonStreamingResponse(resp, c, parsed.ResponseFormat, requestModel)
	}
	if err != nil {
		return nil, err
	}
	if imageCount <= 0 {
		imageCount = parsed.N
	}

	return &OpenAIForwardResult{
		RequestID:    resp.Header.Get("x-request-id"),
		Usage:        usage,
		Model:        requestModel,
		Stream:       parsed.Stream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
		ImageCount:   imageCount,
		ImageSize:    normalizeOpenAIImageSize(parsed.Size),
	}, nil
}
