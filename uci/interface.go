package uci

import "context"

type IEngine interface {
	Options() map[string]Option
	SetOption(name, value string) error
	NewGame(ctx context.Context) error
	Go(ctx context.Context, pos Position, lim Limits, opts GoOpts) (<-chan Info, error)
	Stop(ctx context.Context) error
	Close() error
}
