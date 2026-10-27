package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Addr      string
	DataDir   string
	ConfigDir string
}

func main() {
	if err := run(); err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := parseFlags()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	db, err := store.Open(cfg.DataDir)
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	registry, err := engines.NewRegistry(cfg.ConfigDir)
	if err != nil {
		return fmt.Errorf("initialising engine registry: %w", err)
	}
	defer registry.Close()

	sessionMgr := play.NewSessionManager(db, registry)

	analysisSvc := analysis.NewService(db, registry)

	handler := web.NewHandler(web.Dependencies{
		Store:           db,
		EngineRegistry:  registry,
		SessionManager:  sessionMgr,
		AnalysisService: analysisSvc,
	})

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if !isLoopback(cfg.Addr) {
		slog.Warn("binding to non-loopback address", "addr", cfg.Addr)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		slog.Info("shutting down gracefully")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	analysisSvc.Shutdown(shutdownCtx)

	sessionMgr.Shutdown(shutdownCtx)

	slog.Info("shutdown complete")
	return nil
}

func ParseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.Addr, "addr", getEnvOrDefault("LOCALCHESS_ADDR", "127.0.0.1:8080"),
		"address to listen on")
	flag.StringVar(&cfg.DataDir, "data-dir", getEnvOrDefault("LOCALCHESS_DATA", defaultDataDir()),
		"directory for database and data files")
	flag.StringVar(&cfg.ConfigDir, "config-dir", getEnvOrDefault("LOCALCHESS_CONFIG", defaultConfigDir()),
		"directory for configuration files")

	flag.Parse()
	return cfg
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func defaultDataDir() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return xdg + "/localchess"
	}
	home, _ := os.UserHomeDir()
	return home + "/.local/share/localchess"
}

func defaultConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg + "/localchess"
	}
	home, _ := os.UserHomeDir()
	return home + "/.config/localchess"
}

func isLoopback(addr string) bool {
	// TODO Simple check; production code should parse and validate properly
	return addr == "127.0.0.1:8080" || addr == "localhost:8080" ||
		len(addr) > 0 && (addr[0] == ':' || addr[:9] == "127.0.0.1" || addr[:9] == "localhost")
}
