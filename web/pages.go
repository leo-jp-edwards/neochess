package web

import "net/http"

type PageData struct {
	Title     string
	ActiveNav string
	CSRFToken string
	Board     *BoardView
	Data      any
}

type BoardView struct {
	GameID      string
	FEN         string
	Orientation string
	Turn        string
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderPage(w, "home", "pages/home", PageData{
		Title:     "NeoChess",
		ActiveNav: "home",
		CSRFToken: h.csrf,
	})
}

func (h *Handler) handlePlayPage(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderPage(w, "play", "pages/play", PageData{
		Title:     "Play - NeoChess",
		ActiveNav: "play",
		CSRFToken: h.csrf,
		Board: &BoardView{
			GameID:      "demo",
			FEN:         "start",
			Orientation: "white",
			Turn:        "white",
		},
	})
}

func (h *Handler) handleLibraryPage(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderPage(w, "library", "pages/library", PageData{
		Title:     "Library - NeoChess",
		ActiveNav: "library",
		CSRFToken: h.csrf,
	})
}

func (h *Handler) handleGamePage(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderPage(w, "game", "pages/game", PageData{
		Title:     "Game Page",
		ActiveNav: "game",
		CSRFToken: h.csrf,
	})
}

func (h *Handler) handleListEngines(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderPage(w, "engine", "pages/engine", PageData{
		Title:     "Engine - NeoChess",
		ActiveNav: "engine",
		CSRFToken: h.csrf,
	})
}

func (h *Handler) handleSettingsPage(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderPage(w, "settings", "pages/settings", PageData{
		Title:     "Settings - NeoChess",
		ActiveNav: "settings",
		CSRFToken: h.csrf,
	})
}
