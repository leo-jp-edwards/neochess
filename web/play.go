package web

import "net/http"

func (h *Handler) handleNewGame(w http.ResponseWriter, r *http.Request)      {}
func (h *Handler) handleSelectSquare(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleMakeMove(w http.ResponseWriter, r *http.Request)     {}
func (h *Handler) handleResign(w http.ResponseWriter, r *http.Request)       {}
func (h *Handler) handleEngineSSE(w http.ResponseWriter, r *http.Request)    {}
