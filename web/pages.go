package web

import "net/http"

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request)        {}
func (h *Handler) handlePlayPage(w http.ResponseWriter, r *http.Request)     {}
func (h *Handler) handleLibraryPage(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) handleGamePage(w http.ResponseWriter, r *http.Request)     {}
func (h *Handler) handleSettingsPage(w http.ResponseWriter, r *http.Request) {}
