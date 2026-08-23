package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// opsHandler exposes the OpsService over HTTP with classified error responses.
// Errors are mapped through opsCode/opsHTTPStatus so not_found, conflict,
// transition, policy and invalid are distinguishable from internal failures.
type opsHandler struct {
	service *OpsService
}

func newOpsHandler(service *OpsService) *opsHandler { return &opsHandler{service: service} }

const opsPathPrefix = "/api/v1/operations"

func (h *opsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	opsNoStore(w)
	path := r.URL.Path
	switch {
	case path == opsPathPrefix || path == opsPathPrefix+"/":
		if opsAllowed(r.Method, http.MethodGet) {
			h.search(w, r)
			return
		}
		if opsAllowed(r.Method, http.MethodPost) {
			h.create(w, r)
			return
		}
		opsMethodNotAllowed(w, r)
		return
	case path == opsPathPrefix+"/snapshot":
		h.snapshot(w, r)
		return
	}
	rest := opsPathID(path, opsPathPrefix+"/")
	if rest == "" {
		http.NotFound(w, r)
		return
	}
	parts := strings.SplitN(rest, "/", 2)
	id := parts[0]
	if len(parts) == 1 {
		if opsAllowed(r.Method, http.MethodGet) {
			h.get(w, r, id)
			return
		}
		if opsAllowed(r.Method, http.MethodDelete) {
			h.delete(w, r, id)
			return
		}
		opsMethodNotAllowed(w, r)
		return
	}
	switch parts[1] {
	case "transition":
		if opsAllowed(r.Method, http.MethodPost) {
			h.transition(w, r, id)
			return
		}
		opsMethodNotAllowed(w, r)
		return
	case "audit":
		if opsAllowed(r.Method, http.MethodGet) {
			h.audit(w, r, id)
			return
		}
		opsMethodNotAllowed(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}
}

func (h *opsHandler) create(w http.ResponseWriter, r *http.Request) {
	var record OpsRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		opsWriteError(w, http.StatusBadRequest, "invalid", "invalid JSON body")
		return
	}
	out, err := h.service.Create(r.Context(), record)
	if err != nil {
		opsWriteServiceError(w, err)
		return
	}
	opsJSON(w, http.StatusCreated, out)
}

func (h *opsHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	out, err := h.service.Get(r.Context(), id)
	if err != nil {
		opsWriteServiceError(w, err)
		return
	}
	opsJSON(w, http.StatusOK, out)
}

func (h *opsHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if _, err := h.service.Get(r.Context(), id); err != nil {
		opsWriteServiceError(w, err)
		return
	}
	opsJSON(w, http.StatusNoContent, nil)
}

func (h *opsHandler) transition(w http.ResponseWriter, r *http.Request, id string) {
	var body struct {
		Status   string `json:"status"`
		Revision int    `json:"revision"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		opsWriteError(w, http.StatusBadRequest, "invalid", "invalid JSON body")
		return
	}
	if !opsStatusValid(OpsStatus(body.Status)) {
		opsWriteError(w, http.StatusBadRequest, "invalid", "unsupported status")
		return
	}
	out, err := h.service.Transition(r.Context(), id, body.Revision, OpsStatus(body.Status), opsActorFromRequest(r))
	if err != nil {
		opsWriteServiceError(w, err)
		return
	}
	opsJSON(w, http.StatusOK, out)
}

func (h *opsHandler) audit(w http.ResponseWriter, r *http.Request, id string) {
	opsJSON(w, http.StatusOK, h.service.Audit(id))
}

func (h *opsHandler) search(w http.ResponseWriter, r *http.Request) {
	q := OpsQuery{Page: parseInt(r, "page", 1), PageSize: parseInt(r, "page_size", 25)}
	q.Subject = r.URL.Query().Get("subject")
	q.Status = OpsStatus(r.URL.Query().Get("status"))
	q.Priority = OpsPriority(r.URL.Query().Get("priority"))
	q.Owner = r.URL.Query().Get("owner")
	page, err := h.service.Search(r.Context(), q)
	if err != nil {
		opsWriteServiceError(w, err)
		return
	}
	opsJSON(w, http.StatusOK, page)
}

func (h *opsHandler) snapshot(w http.ResponseWriter, r *http.Request) {
	opsJSON(w, http.StatusOK, h.service.Snapshot())
}

func opsMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "")
	opsWriteError(w, http.StatusMethodNotAllowed, "invalid", "method not allowed")
}

func opsWriteServiceError(w http.ResponseWriter, err error) {
	code := opsCode(err)
	opsWriteError(w, opsHTTPStatus(code), code, err.Error())
}

func opsWriteError(w http.ResponseWriter, status int, code, message string) {
	opsJSON(w, status, map[string]any{"error": code, "message": message})
}

func parseInt(r *http.Request, key string, def int) int {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return def
	}
	n := 0
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return def
		}
		n = n*10 + int(ch-'0')
	}
	if n < 1 {
		return def
	}
	return n
}
