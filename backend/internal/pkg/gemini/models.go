// Package gemini provides minimal fallback model metadata for Gemini native endpoints.
// It is used when upstream model listing is unavailable (e.g. OAuth token missing AI Studio scopes).
package gemini

type Model struct {
	Name                       string   `json:"name"`
	DisplayName                string   `json:"displayName,omitempty"`
	Description                string   `json:"description,omitempty"`
	SupportedGenerationMethods []string `json:"supportedGenerationMethods,omitempty"`
}

type ModelsListResponse struct {
	Models []Model `json:"models"`
}

func DefaultModels() []Model {
	methods := []string{"generateContent", "streamGenerateContent"}
	return []Model{
		{Name: "models/gemini-2.5-flash", SupportedGenerationMethods: methods},
		{Name: "models/gemini-2.5-flash-lite", SupportedGenerationMethods: methods},
		{Name: "models/gemini-2.5-pro", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3-flash", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3-flash-preview", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3-pro", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3-pro-high", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3-pro-preview", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-pro", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-pro-preview", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-3x2", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2x3", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-3x4", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4x3", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4x5", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-5x4", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-9x16", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-16x9", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-21x9", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k-3x2", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k-2x3", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k-3x4", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k-4x3", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k-4x5", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k-5x4", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k-9x16", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k-16x9", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-2k-21x9", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k-3x2", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k-2x3", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k-3x4", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k-4x3", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k-4x5", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k-5x4", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k-9x16", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k-16x9", SupportedGenerationMethods: methods},
		{Name: "models/gemini-3.1-flash-image-4k-21x9", SupportedGenerationMethods: methods},
	}
}

func FallbackModelsList() ModelsListResponse {
	return ModelsListResponse{Models: DefaultModels()}
}

func FallbackModel(model string) Model {
	methods := []string{"generateContent", "streamGenerateContent"}
	if model == "" {
		return Model{Name: "models/unknown", SupportedGenerationMethods: methods}
	}
	if len(model) >= 7 && model[:7] == "models/" {
		return Model{Name: model, SupportedGenerationMethods: methods}
	}
	return Model{Name: "models/" + model, SupportedGenerationMethods: methods}
}
