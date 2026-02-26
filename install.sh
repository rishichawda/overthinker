#!/usr/bin/env bash
set -euo pipefail

REPO="rishichawda/overthinker"
BINARY="overthink"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# ── helpers ──────────────────────────────────────────────────────────────────

info()  { printf "\033[1;34m[overthink]\033[0m %s\n" "$*"; }
error() { printf "\033[1;31m[overthink error]\033[0m %s\n" "$*" >&2; exit 1; }

need() {
  command -v "$1" &>/dev/null || error "Required tool not found: $1"
}

# ── detect OS and arch ───────────────────────────────────────────────────────

detect_os() {
  case "$(uname -s)" in
    Darwin) echo "Darwin" ;;
    Linux)  echo "Linux"  ;;
    *)      error "Unsupported OS: $(uname -s)" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64)  echo "x86_64" ;;
    arm64|aarch64) echo "arm64"  ;;
    *)             error "Unsupported architecture: $(uname -m)" ;;
  esac
}

# ── resolve latest version ───────────────────────────────────────────────────

resolve_version() {
  if [[ -n "${VERSION:-}" ]]; then
    echo "$VERSION"
    return
  fi
  need curl
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed -E 's/.*"v?([^"]+)".*/\1/'
}

# ── checksum verification ────────────────────────────────────────────────────

# Prefer sha256sum (standard on Linux); fall back to shasum (macOS / Perl)
sha256_verify() {
  if command -v sha256sum &>/dev/null; then
    sha256sum -c --status
  else
    shasum -a 256 -c --status
  fi
}

# ── main ─────────────────────────────────────────────────────────────────────

main() {
  need curl
  need tar

  local os arch version tarball url checksum_url tmp_dir

  os="$(detect_os)"
  arch="$(detect_arch)"
  version="$(resolve_version)"

  info "Installing ${BINARY} v${version} (${os}/${arch})"

  tarball="${BINARY}_${os}_${arch}.tar.gz"
  url="https://github.com/${REPO}/releases/download/v${version}/${tarball}"
  checksum_url="https://github.com/${REPO}/releases/download/v${version}/checksums.txt"

  tmp_dir="$(mktemp -d)"
  trap 'rm -rf "$tmp_dir"' EXIT

  info "Downloading ${tarball}..."
  curl -fsSL "$url" -o "${tmp_dir}/${tarball}"

  info "Verifying checksum..."
  curl -fsSL "$checksum_url" -o "${tmp_dir}/checksums.txt"
  grep "${tarball}" "${tmp_dir}/checksums.txt" \
    | (cd "$tmp_dir" && sha256_verify) \
    || error "Checksum verification failed. The download may be corrupted."

  info "Extracting..."
  tar -xzf "${tmp_dir}/${tarball}" -C "$tmp_dir"

  if [[ -w "$INSTALL_DIR" ]] || [[ "$(id -u)" -eq 0 ]]; then
    install -m 755 "${tmp_dir}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    info "Installed to ${INSTALL_DIR}/${BINARY}"
  else
    local user_dir="${HOME}/.local/bin"
    mkdir -p "$user_dir"
    install -m 755 "${tmp_dir}/${BINARY}" "${user_dir}/${BINARY}"
    info "Installed to ${user_dir}/${BINARY}"
    info "Make sure ${user_dir} is in your PATH:"
    info "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc"
  fi

  info "Done! Run: ${BINARY} \"Should I use a bash install script?\""
}

main "$@"
