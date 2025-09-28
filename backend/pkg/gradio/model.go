package gradio

type dataFile struct {
	Meta     map[string]any `json:"meta,omitempty"`
	MimeType string         `json:"mime_type,omitempty"`
	Path     string         `json:"path"`
	OrigName string         `json:"orig_name,omitempty"`
	Size     int64          `json:"size,omitempty"`
	URL      string         `json:"url,omitempty"`
}

func newDataFile(mimeType string, path string, origName string, size int64, url string) *dataFile {
	return &dataFile{
		Meta: map[string]any{
			"_type": "gradio.FileData",
		},
		MimeType: mimeType,
		Path:     path,
		OrigName: origName,
		Size:     size,
		URL:      url,
	}
}

type bodyPredict struct {
	Data []any `json:"data"`
}

type respPredict struct {
	EventID string `json:"event_id"`
}

type respJWT struct {
	Token       string `json:"token"`
	AccessToken string `json:"accessToken"`
	Expire      int64  `json:"exp"`
}
