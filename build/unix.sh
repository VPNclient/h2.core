#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

OUT_DIR="${OUT_DIR:-dist}"
NAME="${NAME:-https_vpn}"
VERSION="${VERSION:-}"
TAGS="${TAGS:-}"

if [[ -z "${VERSION}" ]]; then
  VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo "0.1.0-dev")"
fi

mkdir -p "$OUT_DIR"

build_one() {
  local goos="$1"
  local goarch="$2"
  local outfile="$OUT_DIR/${NAME}_${goos}_${goarch}"
  if [[ "$goos" == "darwin" ]]; then
    outfile="$OUT_DIR/${NAME}"
  fi

  echo "Building $outfile (GOOS=$goos GOARCH=$goarch TAGS=${TAGS:-<none>})"
  env CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -ldflags "-s -w -X main.Version=$VERSION" \
    ${TAGS:+-tags "$TAGS"} \
    -o "$outfile" ./cmd/https-vpn
}

# Native-friendly defaults: build darwin/<host arch> as dist/https_vpn,
# and cross-compile linux binaries for server usage.
HOST_ARCH="$(go env GOARCH)"
build_one darwin "$HOST_ARCH"
build_one linux amd64
build_one linux arm64

echo "Done. Outputs in $OUT_DIR/"
