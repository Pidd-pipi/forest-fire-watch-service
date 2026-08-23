package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func opsEnterpriseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := newOpsResponseRecorder(w)
		rec.Header().Set("X-Operations-Domain", opsDomainName)
		if strings.TrimSpace(r.Header.Get("X-Request-ID")) == "" {
			rec.Header().Set("X-Operations-Request", "generated")
		} else {
			rec.Header().Set("X-Operations-Request", "provided")
		}
		next.ServeHTTP(rec, r)
		// Measure once the handler has finished, then flush headers so the
		// latency marker is guaranteed to travel with the response.
		rec.Header().Set("X-Operations-Latency-Ms", formatOpsInt(int(time.Since(start).Milliseconds())))
		rec.flush()
	})
}

// opsResponseRecorder buffers status and body so trailing middleware work
// (such as setting the latency header in a defer after the handler returns)
// can still attach response headers before they hit the wire.
type opsResponseRecorder struct {
	header      http.Header
	wroteHeader bool
	status      int
	buf         bytes.Buffer
	w           http.ResponseWriter
}

func newOpsResponseRecorder(w http.ResponseWriter) *opsResponseRecorder {
	return &opsResponseRecorder{header: http.Header{}, w: w, status: 0}
}
func (r *opsResponseRecorder) Header() http.Header { return r.header }
func (r *opsResponseRecorder) WriteHeader(code int) {
	if r.wroteHeader {
		return
	}
	r.status = code
	r.wroteHeader = true
}
func (r *opsResponseRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.buf.Write(b)
}
func (r *opsResponseRecorder) flush() {
	dst := r.w.Header()
	for k, vs := range r.header {
		dst[k] = append([]string(nil), vs...)
	}
	if !r.wroteHeader {
		r.status = http.StatusOK
	}
	r.w.WriteHeader(r.status)
	r.w.Write(r.buf.Bytes())
}

func formatOpsInt(value int) string {
	if value == 0 {
		return "0"
	}
	out := ""
	for value > 0 {
		out = string(rune('0'+value%10)) + out
		value /= 10
	}
	return out
}
func opsJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func opsAllowed(method string, allowed ...string) bool {
	for _, candidate := range allowed {
		if method == candidate {
			return true
		}
	}
	return false
}
func opsPathID(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.Trim(strings.TrimPrefix(path, prefix), "/")
}
func opsActorFromRequest(r *http.Request) string {
	value := strings.TrimSpace(r.Header.Get("X-Operator"))
	if value == "" {
		return "web"
	}
	return value
}
func opsNoStore(w http.ResponseWriter)    { w.Header().Set("Cache-Control", "no-store") }
func opsRequestID(r *http.Request) string { return r.Header.Get("X-Request-ID") }
