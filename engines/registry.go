package engines

import (
	"context"
	"fmt"
	"neochess/uci"
	"os"
	"path/filepath"
)

// EngineConfig represents a configured UCI engine
type EngineConfig struct {
	ID      string            `toml:"id"`
	Name    string            `toml:"name"`
	Path    string            `toml:"path"`
	Args    []string          `toml:"args"`
	Options map[string]string `toml:"options"`
}

// CachedOptions holds probed engine capabilities
type CachedOptions struct {
	EngineID        string
	Options         map[string]uci.Option
	SupportsElo     bool
	EloMin          int
	EloMax          int
	SupportsMultiPV bool
	ProbedAt        int64
}

// Registry managed configured engines and their cached options
type Registry struct {
	configPath string
	engines    map[string]*EngineConfig
	cache      map[string]*CachedOptions
}

// NewRegistry loads engin configuration from the config directory
func NewRegistry(configDir string) (*Registry, error) {
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return nil, fmt.Errorf("engines: creating config dir: %w", err)
	}

	configPath := filepath.Join(configDir, "engines.toml")

	r := &Registry{
		configPath: configPath,
		engines:    make(map[string]*EngineConfig),
		cache:      make(map[string]*CachedOptions),
	}

	if err := r.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("engines: loading config: %w", err)
	}

	return r, nil
}

// load reads the engines.toml config file
func (r *Registry) load() error {
	// TODO: parse toml config file
	return nil
}

// Save persists the currrent engine configuration
func (r *Registry) Save() error {
	// TODO: write toml config file
	return nil
}

// List returns all configured engines
func (r *Registry) List() []*EngineConfig {
	result := make([]*EngineConfig, 0, len(r.engines))
	for _, e := range r.engines {
		result = append(result, e)
	}
	return result
}

// Get retrieves an engine by ID
func (r *Registry) Get(id string) (*EngineConfig, bool) {
	e, ok := r.engines[id]
	return e, ok
}

// Add registers a new engine, probing it for options
func (r *Registry) Add(ctx context.Context, cfg *EngineConfig) error {
	// TODO start engine, run handshake, cache options and add to registry
	return nil
}

// Remove deletes an engine from the registry
func (r *Registry) Remove(id string) error {
	delete(r.engines, id)
	delete(r.cache, id)
	return r.Save()
}

// Detect finds common engines on PATH.
func (r *Registry) Detect(ctx context.Context) ([]*EngineConfig, error) {
	// TODO: Look for stockfish, lc0, etc. on PATH
	// Return detected but not yet saved engines
	return nil, nil
}

// GetCachedOptions returns cached options for an engine.
func (r *Registry) GetCachedOptions(id string) (*CachedOptions, bool) {
	c, ok := r.cache[id]
	return c, ok
}

// StartEngine launches an engine process ready for use
func (r *Registry) StartEngine(ctx context.Context, id string) (*uci.Engine, error) {
	cfg, ok := r.engines[id]
	if !ok {
		return nil, fmt.Errorf("engines: unknown engine %q", id)
	}

	engine, err := uci.Start(ctx, cfg.Path, cfg.Args)
	if err != nil {
		return nil, fmt.Errorf("engines: starting %s: %w", id, err)
	}

	for name, value := range cfg.Options {
		if err := engine.SetOption(name, value); err != nil {
			engine.Close()
			return nil, fmt.Errorf("engines: setting option %s: %w", name, err)
		}
	}

	return engine, nil
}

// Close cleans up any resources held by the registry
func (r *Registry) Close() error {
	return nil
}
