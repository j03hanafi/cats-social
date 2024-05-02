package handler

type baseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type version struct {
	Version string `json:"version"`
}
