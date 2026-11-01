#!/usr/bin/env bash
set -euo pipefail

HTMX_VERSION="${HTMX_VERSION:-2.0.4}"
CHESSGROUND_VERSION="${CHESSGROUND_VERSION:-10.1.1}"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STATIC_DIR="$ROOT/web/static"

mkdir -p "$STATIC_DIR/js" "$STATIC_DIR/css" "$STATIC_DIR/vendor/licenses"

echo "Vendoring htmx ${HTMX_VERSION}"
curl -fsSL \
  "https://unpkg.com/htmx.org@${HTMX_VERSION}/dist/htmx.min.js" \
  -o "$STATIC_DIR/js/htmx.min.js"

echo "Vendoring htmx SSE extension"
curl -fsSL \
  "https://unpkg.com/htmx-ext-sse@2.2.2/sse.js" \
  -o "$STATIC_DIR/js/htmx-sse.js"

echo "Vendoring chessground ${CHESSGROUND_VERSION}"
curl -fsSL \
  "https://cdn.jsdelivr.net/npm/@lichess-org/chessground@${CHESSGROUND_VERSION}/dist/chessground.min.js" \
  -o "$STATIC_DIR/js/chessground.min.js"

cat > "$STATIC_DIR/vendor/manifest.txt" <<EOF
htmx.org ${HTMX_VERSION}
htmx-ext-sse 2.2.2
@lichess-org/chessground ${CHESSGROUND_VERSION}
EOF

echo "Done."
