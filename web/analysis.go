package web

import "net/http"

func (h *Handler) handleStartAnalysis(w http.ResponseWriter, r *http.Request)     {}
func (h *Handler) handleCancelAnalysis(w http.ResponseWriter, r *http.Request)    {}
func (h *Handler) handleAnalysisSSE(w http.ResponseWriter, r *http.Request)       {}
func (h *Handler) handleStartLiveAnalysis(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleStopLiveAnalysis(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) handleLiveAnalysisSSE(w http.ResponseWriter, r *http.Request)   {}
