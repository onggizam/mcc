#!/usr/bin/env bash
set -euo pipefail

BUMP="${BUMP:-auto}"

LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

RANGE=""
if git rev-parse -q --verify "refs/tags/${LAST_TAG}" >/dev/null 2>&1; then
  RANGE="${LAST_TAG}..HEAD"
else
  RANGE="HEAD"
fi

if [[ "$BUMP" == "auto" ]]; then
  COMMITS=$(git log --format=%B $RANGE || true)

  shopt -s nocasematch
  if echo "$COMMITS" | grep -Eq '^.*!:' || echo "$COMMITS" | grep -Eq 'BREAKING CHANGE:'; then
    BUMP="major"
  elif echo "$COMMITS" | grep -Eq '^feat(\(.*\))?:'; then
    BUMP="minor"
  elif echo "$COMMITS" | grep -Eq '^(fix|perf|refactor)(\(.*\))?:'; then
    BUMP="patch"
  else
    BUMP="patch"
  fi
  shopt -u nocasematch
fi

clean="${LAST_TAG#v}"
IFS='.' read -r MA MI PA <<<"$clean"

inc() {
  local bump="$1"
  local ma="$2" mi="$3" pa="$4"
  case "$bump" in
    major) ma=$((ma+1)); mi=0; pa=0 ;;
    minor) mi=$((mi+1)); pa=0 ;;
    patch) pa=$((pa+1)) ;;
    *) echo "unknown bump: $bump" >&2; exit 1 ;;
  esac
  echo "v${ma}.${mi}.${pa}"
}

NEXT_TAG=$(inc "$BUMP" "$MA" "$MI" "$PA")

if [[ -d pkg/version ]]; then
  cat > pkg/version/version.go <<EOF
package version

const Version = "${NEXT_TAG}"
EOF
fi

echo "${NEXT_TAG}"