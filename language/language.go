package language

const (
	English    string = "en"
	Indonesian string = "id"
)

func HTTPStatusText(lang string, code int) string {
	var statusText string

	switch lang {
	case English:
		statusText = statusTextEn[code]
	case Indonesian:
		statusText = statusTextId[code]
	}

	return statusText
}
