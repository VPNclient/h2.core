#!/usr/bin/env bash
#
# Build h2.core as a C-compatible shared library (.so/.dylib)
#
# Usage:
#   ./build/cgo.sh              # Build for current platform
#   ./build/cgo.sh all          # Cross-compile for all platforms
#
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

OUT_DIR="${OUT_DIR:-dist}"
VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "0.1.0-dev")}"

mkdir -p "$OUT_DIR"

build_shared() {
    local goos="$1"
    local goarch="$2"
    local ext="$3"
    local outfile="$OUT_DIR/libh2core_${goos}_${goarch}${ext}"

    echo "Building $outfile (GOOS=$goos GOARCH=$goarch)"

    # CGO must be enabled for c-shared
    env CGO_ENABLED=1 GOOS="$goos" GOARCH="$goarch" \
        go build -buildmode=c-shared \
        -ldflags "-s -w -X main.Version=$VERSION" \
        -o "$outfile" ./cgo

    # The .h file is generated automatically, rename it
    if [[ -f "${outfile%.${ext}}.h" ]]; then
        mv "${outfile%.${ext}}.h" "$OUT_DIR/h2core.h"
    fi

    echo "  -> $outfile"
}

# Detect host platform
HOST_OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
HOST_ARCH="$(go env GOARCH)"

case "$HOST_OS" in
    darwin)
        HOST_OS="darwin"
        EXT=".dylib"
        ;;
    linux)
        HOST_OS="linux"
        EXT=".so"
        ;;
    *)
        echo "Unsupported OS: $HOST_OS"
        exit 1
        ;;
esac

if [[ "${1:-}" == "all" ]]; then
    echo "Cross-compiling for all platforms..."
    echo "Note: Cross-compilation requires appropriate C cross-compilers"
    echo ""

    # Native build (always works)
    build_shared "$HOST_OS" "$HOST_ARCH" "$EXT"

    # Cross-compile warnings
    if [[ "$HOST_OS" == "darwin" ]]; then
        echo ""
        echo "For Linux cross-compilation on macOS, install:"
        echo "  brew install FiloSottile/musl-cross/musl-cross"
        echo ""
        # Attempt Linux builds if cross-compiler available
        if command -v x86_64-linux-musl-gcc &> /dev/null; then
            CC=x86_64-linux-musl-gcc build_shared linux amd64 .so
        fi
        if command -v aarch64-linux-musl-gcc &> /dev/null; then
            CC=aarch64-linux-musl-gcc build_shared linux arm64 .so
        fi
    elif [[ "$HOST_OS" == "linux" ]]; then
        # Linux native builds
        build_shared linux amd64 .so
        if [[ "$HOST_ARCH" == "arm64" ]] || command -v aarch64-linux-gnu-gcc &> /dev/null; then
            build_shared linux arm64 .so
        fi
    fi
else
    # Build for current platform only
    build_shared "$HOST_OS" "$HOST_ARCH" "$EXT"
fi

# Copy header file
cp cgo/h2core.h "$OUT_DIR/"

echo ""
echo "Build complete. Outputs in $OUT_DIR/"
ls -la "$OUT_DIR"/libh2core* "$OUT_DIR"/h2core.h 2>/dev/null || true
