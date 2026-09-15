#!/bin/bash
set -e
cd "$(dirname "$0")/.."

BUILD_OUTPUT="zincsearch"
MAIN_PACKAGE="cmd/zincsearch/main.go"
FRONTEND_DIR="frontend"

echo "Cleaning previous builds..."
rm -f "$BUILD_OUTPUT"

echo "Building frontend assets..."
(cd "$FRONTEND_DIR" && npm ci && npm run build)

# Get build metadata
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
BUILD_DATE=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT_HASH=$(git rev-parse HEAD 2>/dev/null || echo "unknown")

# Detect platform and architecture
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|x64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    i386|i686) ARCH="386" ;;
esac

case "$PLATFORM" in
    linux)
        GOOS="linux"
        ;;
    darwin)
        GOOS="darwin"
        ;;
    *)
        echo "Warning: Unknown platform '$PLATFORM', defaulting to Linux"
        GOOS="linux"
        ;;
esac

if [[ "$ARCH" != "amd64" && "$ARCH" != "arm64" ]]; then
    echo "Warning: Architecture '$ARCH' may not be fully supported, building anyway..."
fi

echo "Building for $GOOS/$ARCH..."
echo "Version: $VERSION"
echo "Build date: $BUILD_DATE"
echo "Commit: $COMMIT_HASH"

# Build command with common flags
BUILD_FLAGS=(
    -ldflags="-w -s \
        -X github.com/zincsearch/zincsearch/pkg/meta.Version=${VERSION} \
        -X github.com/zincsearch/zincsearch/pkg/meta.BuildDate=${BUILD_DATE} \
        -X github.com/zincsearch/zincsearch/pkg/meta.CommitHash=${COMMIT_HASH}"
    -trimpath
    -tags ne  # drop gse embedded dictionaries, same as .goreleaser.yml
    -o "$BUILD_OUTPUT"
)

CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$ARCH" go build "${BUILD_FLAGS[@]}" "$MAIN_PACKAGE"

chmod +x "$BUILD_OUTPUT"
if [ -f "$BUILD_OUTPUT" ]; then
    echo "✅ Build completed successfully: $BUILD_OUTPUT"
    file "$BUILD_OUTPUT"
    ls -lh "$BUILD_OUTPUT"
else
    echo "❌ Build failed: $BUILD_OUTPUT not found"
    exit 1
fi