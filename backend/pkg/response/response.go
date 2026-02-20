package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Data    any      `json:"data,omitempty"`
	Error   string   `json:"error,omitempty"`
	Code    string   `json:"code,omitempty"`
	TraceID string   `json:"trace_id,omitempty"`
	Meta    MetaInfo `json:"meta,omitempty"`
}

type MetaInfo struct {
	Disclosures []string `json:"disclosures,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(200, Envelope{Data: data, TraceID: traceID(c)})
}

func Created(c *gin.Context, data any) {
	c.JSON(201, Envelope{Data: data, TraceID: traceID(c)})
}

func BadRequest(c *gin.Context, msg string) {
	fail(c, 400, "bad_request", msg)
}

func NotFound(c *gin.Context, msg string) {
	fail(c, 404, "not_found", msg)
}

func Internal(c *gin.Context, msg string) {
	fail(c, 500, "internal_error", msg)
}

func TooManyRequests(c *gin.Context, msg string) {
	fail(c, 429, "rate_limited", msg)
}

func fail(c *gin.Context, status int, code string, msg string) {
	c.JSON(status, Envelope{Error: msg, Code: code, TraceID: traceID(c)})
}

func traceID(c *gin.Context) string {
	v, ok := c.Get("request_id")
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
