package requests

type BatchRequest struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   string `json:"body"`
}
