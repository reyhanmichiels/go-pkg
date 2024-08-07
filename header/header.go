package header

const (
	// Header keys
	KeyRequestID      string = "x-request-id"
	KeyAuthorization  string = "authorization"
	KeyUserAgent      string = "user-agent"
	KeyContentType    string = "content-type"
	KeyContentAccept  string = "accept"
	KeyAcceptLanguage string = "accept-language"
	KeyCacheControl   string = "cache-control"
	KeyDeviceType     string = "x-device-type"
	KeyServiceName    string = "x-service-name"
	KeyEventType      string = "x-event-type"
	KeyEventSource    string = "x-event-source"
	KeyDeviceToken    string = "x-device-token"

	// Content type. Specifying the payload in the request
	ContentTypeJSON string = "application/json"
	ContentTypeXML  string = "application/xml"
	ContentTypeForm string = "application/x-www-form-urlencoded"

	// Accepting media. Specifying the types of requested media (in the response)
	// See here: https://en.wikipedia.org/wiki/Content_negotiation
	MediaTextPlain string = "text/plain"
	MediaTextHTML  string = "text/html"
	MediaTextCSV   string = "text/csv"
	MediaTextXML   string = "text/xml"

	MediaImageGIF  string = "image/gif"
	MediaImageJPEG string = "image/jpeg"
	MediaImagePNG  string = "image/png"
	MediaImageWEBP string = "image/webp"

	// Cache control
	CacheControlNoCache string = "no-cache"
	CacheControlNoStore string = "no-store"
)
