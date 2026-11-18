package queue

type MediaMetadata struct {
	DurationSeconds float64 `json:"duration_seconds"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	Codec           string  `json:"codec"`
}

func ValidateMetadata(m *MediaMetadata) bool {
	return m != nil && m.DurationSeconds > 0 && m.Width > 0 && m.Height > 0
}
