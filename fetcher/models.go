package fetcher

type UrlFetchRequest struct {
	ExecutionTimeout *int         `json:"execution_timeout"`
	UrlRequests      []UrlRequest `json:"url_requests" validate:"required,dive"`
}

type UrlRequest struct {
	Url     string            `json:"url" validate:"required,url"`
	Timeout *int              `json:"timeout"`
	Headers map[string]string `json:"headers"`
}

type UrlFetchResult struct {
	Url     string `json:"url"`
	Code    int    `json:"code"`
	Payload string `json:"payload,omitempty"`
	Error   string `json:"error,omitempty"`
}

type UrlFetchResponse struct {
	Results []UrlFetchResult `json:"results"`
}
