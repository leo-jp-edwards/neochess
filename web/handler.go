// Package web provides HTTP handlers, templates and SSE support

package web

import (
	"embed"
	"neochess/analysis"
	"neochess/engines"
	"neochess/play"
	"neochess/store"
	"net/http"
)

//go:embed templates/*.html static/*
var context embed.FS

// Dependencies holds all services needed by the web handlers
type Dependencies struct {
	Store           store.Store
	EngineRegistry  *engines.Registry
	SessionManager  *play.SessionManager
	AnalysisService *analysis.Service
}

// Handler is the main HTTP handler for the application
type Handler struct {
	deps Dependencies
	mux  *http.ServeMux
	csrf string // per-process CSRF token
}

// NewHandler creates a new HTTP handler with all routes configured
func NewHandler(deps Dependencies) *Handler {
	h := &Handler{
		deps: deps,
		mux:  http.NewServeMux(),
		csrf: generateCSRFToken(),
	}

	h.setupRoutes()
	return h
}

func (h *Handler) setupRoutes() {
	// Static assets
	h.mux.HandleFunc("GET /static/", h.handleStatic)

	// Pages
	h.mux.HandleFunc("GET /", h.handleIndex)
	h.mux.HandleFunc("GET /play", h.handlePlayPage)
	h.mux.HandleFunc("GET /library", h.handleLibraryPage)
	h.mux.HandleFunc("GET/library/{id}", h.handleGamePage)
	h.mux.HandleFunc("GET /settings", h.handleSettingsPage)

	// Play API
	h.mux.HandleFunc("POST /play/new", h.handleNewGame)
	h.mux.HandleFunc("POST /play/{id}/select", h.handleSelectSquare)
	h.mux.HandleFunc("POST /play/{id}/move", h.handleMakeMove)
	h.mux.HandleFunc("POST /play/{id}/resign", h.handleResign)
	h.mux.HandleFunc("GET /play/{id}/engine-sse", h.handleEngineSSE)

	// Library API
	h.mux.HandleFunc("POST /library/import", h.handleImport)
	h.mux.HandleFunc("GET /library/export", h.handleExport)
	h.mux.HandleFunc("DELETE /library/{id}", h.handleDeleteGame)

	// Analysis API
	h.mux.HandleFunc("POST /analysis/{id}/start", h.handleStartAnalysis)
	h.mux.HandleFunc("POST /analysis/{id}/cancel", h.handleCancelAnalysis)
	h.mux.HandleFunc("GET /analysis/{id}/progress-sse", h.handleAnalysisSSE)
	h.mux.HandleFunc("POST /analysis/live/start", h.handleStartLiveAnalysis)
	h.mux.HandleFunc("POST /analysis/live/stop", h.handleStopLiveAnalysis)
	h.mux.HandleFunc("GET /analysis/live/sse", h.handleLiveAnalysisSSE)

	// Engine management
	h.mux.HandleFunc("GET /engines", h.handleListEngines)
	h.mux.HandleFunc("POST /engines", h.handleAddEngine)
	h.mux.HandleFunc("DELETE /engines/{id}", h.handleRemoveEngine)
	h.mux.HandleFunc("POST /engines/detect", h.handleDetectEngines)
}

// ServeHTTP implements http.Handler with security middleware
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Security middleware
	if !h.checkHost(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	if !h.checkOrigin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	h.mux.ServeHTTP(w, r)
}

func (h *Handler) checkHost(r *http.Request) bool {
	// TODO Verify Host Header is loopback
	return true
}

func (h *Handler) checkOrigin(r *http.Request) bool {
	// TODO For state-changing requests, verify Origin/Referer matches host
	return true
}

func generateCSRFToken() string {
	// TODO Generate random token
	return "csrf-token-placeholder"
}

// Stub handlers - each would render templates or return htmx fragments

func (h *Handler) handleStatic(w http.ResponseWriter, r *http.Request) {}
