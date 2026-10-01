#!/usr/bin/env bash
# Copyright 2026 eat-pray-ai & OpenWaygate
# SPDX-License-Identifier: Apache-2.0

set -Eeuo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

# ---------------------------------------------------------------------------
# Constants & Helpers
# ---------------------------------------------------------------------------
SENSITIVE_FILES=("client_secret.json" "youtube.token.json" ".env" ".env.local")
RESOURCE_TARGETS=("main.go" "internal/tools/skillgen/main.go" "internal/tools/cmdtestgen/main.go")
DOC_FILES=("README.md" "README_zh.md" "docs/FEATURES.md" "server.json")
ADDLICENSE_ARGS=(
    -c "eat-pray-ai & OpenWaygate"
    -f LICENSE
    -s=only
    -y "$(date +%Y)"
)

is_gitmoji() { [[ "$1" =~ ^:[a-z0-9_]+:[[:space:]]+.+ ]]; }
info()    { printf "\033[1;34m[INFO]\033[0m %s\n" "$*"; }
success() { printf "\033[1;32m[PASS]\033[0m %s\n" "$*"; }
warn()    { printf "\033[1;33m[WARN]\033[0m %s\n" "$*"; }
error()   { printf "\033[1;31m[FAIL]\033[0m %s\n" "$*"; }

abort() {
    error "$*"
    exit 1
}

TMP_FILES=()
cleanup() {
    if [[ ${#TMP_FILES[@]} -gt 0 ]]; then
        rm -f "${TMP_FILES[@]}" 2>/dev/null || true
    fi
}
trap cleanup EXIT INT TERM
trap 'error "Hook aborted unexpectedly at line $LINENO in command: $BASH_COMMAND"' ERR

run_quiet() {
    local tmp_output
    tmp_output="$(mktemp)"
    if ! "$@" > "$tmp_output" 2>&1; then
        error "Command failed: $*"
        abort "See log: $tmp_output"
    fi
    TMP_FILES+=("$tmp_output")
}

run_step() {
    local label="$1"
    shift
    info "Running $label..."
    run_quiet "$@"
    success "$label passed."
}

# ---------------------------------------------------------------------------
# Tier 1: pre-commit
# Ref: docs/BEFORE_RELEASE.md (Section 1: Code & Resource Registration, License, Lint)
# Less overlap: Fast code hygiene and generation checks. Does not run heavy tests.
# ---------------------------------------------------------------------------
hook_pre_commit() {
    info "Tier 1: Running pre-commit hygiene & code generation checks..."

    local changed_files
    changed_files="$(git diff --cached --name-only)"
    if [[ -z "$changed_files" ]]; then
        changed_files="$(git diff --name-only HEAD 2>/dev/null || git status --porcelain | awk '{print $NF}')"
    fi

    info "Checking for sensitive files in git staging area..."
    for file in "${SENSITIVE_FILES[@]}"; do
        if echo "$changed_files" | grep -qx "$file"; then
            abort "Attempting to stage sensitive file: '$file'. Unstage it before committing!"
        fi
    done
    success "No sensitive files detected in staging area."

    if echo "$changed_files" | grep -qE '^(pkg/|cmd/|main\.go|internal/tools/)'; then
        info "Verifying resource registration across main.go and code generators..."
        for dir in pkg/*/; do
            pkg_name="$(basename "$dir")"
            if [[ "$pkg_name" == "auth" || "$pkg_name" == "common" || "$pkg_name" == "utils" ]]; then
                continue
            fi

            import_path="github.com/eat-pray-ai/yutu/cmd/${pkg_name}"
            for target in "${RESOURCE_TARGETS[@]}"; do
                if ! grep -q "$import_path" "$target"; then
                    abort "Resource '$pkg_name' is missing registration '$import_path' in $target (see docs/BEFORE_RELEASE.md Section 1)."
                fi
            done
        done
        success "All resource packages are properly registered."

        info "Running code generators (skillgen & cmdtestgen)..."
        go run ./internal/tools/skillgen
        go run ./internal/tools/cmdtestgen

        if ! git diff --quiet skills/ cmd/agent/ scripts/command-test.sh; then
            abort "Code generators produced unstaged changes. Please stage them before committing."
        fi
        success "Generated skills and smoke tests are up to date."
    fi

    if command -v addlicense >/dev/null 2>&1; then
        local files_to_license
        files_to_license=$(echo "$changed_files" | grep -vE '(\.json$|\.md$|\.yml$|\.yaml$|\.bazel$|^\.idea/)' || true)
        if [[ -n "$files_to_license" ]]; then
            info "Applying license headers to changed files..."
            echo "$files_to_license" | xargs addlicense "${ADDLICENSE_ARGS[@]}"
            success "License headers applied."
        fi
    else
        warn "'addlicense' not found in PATH; skipping license header check."
    fi

    if echo "$changed_files" | grep -qE '(\.go$|^go\.(mod|sum)$)'; then
        info "Checking go mod tidy..."
        go mod tidy
        if ! git diff --exit-code go.mod go.sum >/dev/null 2>&1; then
            abort "go.mod or go.sum has uncommitted modifications after 'go mod tidy'."
        fi
        success "Go module dependencies are clean and tidy."

        if command -v golangci-lint >/dev/null 2>&1; then
            run_step "golangci-lint" golangci-lint run
        else
            warn "'golangci-lint' not found in PATH; running 'go vet ./...' as fallback."
            run_step "go vet" go vet ./...
        fi
    fi
}

# ---------------------------------------------------------------------------
# Tier 2: commit-msg
# Ref: docs/BEFORE_RELEASE.md (Section 4: Gitmoji Commit Messages)
# Less overlap: Only validates the commit message formatting.
# ---------------------------------------------------------------------------
hook_commit_msg() {
    info "Tier 2: Validating commit message (Gitmoji format)..."

    local msg="${1:-}"

    if [[ -f "$msg" ]]; then
        # Strip comments and blank lines to get the actual commit subject
        msg="$(grep -v '^[[:space:]]*#' "$msg" | grep -v '^[[:space:]]*$' | head -n 1 || true)"
    elif [[ -z "$msg" ]] && git rev-parse --verify HEAD >/dev/null 2>&1; then
        msg="$(git log -1 --pretty=%B | head -n 1)"
    fi

    if [[ -z "$msg" ]]; then
        abort "No commit message provided to validate."
    fi

    if ! is_gitmoji "$msg"; then
        error "Invalid commit message: '$msg'"
        echo ""
        echo "Commit messages must follow the Gitmoji convention:"
        echo "  :<gitmoji>: <Summary>"
        echo ""
        echo "Examples:"
        echo "  :sparkles: Add new playlist tool"
        echo "  :bug: Fix token verifier scope check"
        echo "  :memo: Update pre-release documentation"
        echo "  :recycle: Refactor common provider"
        echo ""
        echo "See docs/BEFORE_RELEASE.md Section 4 and https://gitmoji.dev for details."
        exit 1
    fi

    success "Commit message conforms to Gitmoji convention: '$msg'"
}

# ---------------------------------------------------------------------------
# Tier 3: pre-push
# Dispatches to hook_pre_push_tag for tag pushes, or hook_pre_push_branch otherwise.
# ---------------------------------------------------------------------------
hook_pre_push() {
    local target="${1:-}"
    local pushed_tag=""

    if [[ "$target" =~ ^v[0-9]+ ]]; then
        pushed_tag="$target"
    elif [[ "$target" == "tag" ]]; then
        pushed_tag="${2:-}"
    elif [[ ! -t 0 ]]; then
        while read -r local_ref _; do
            [[ -z "$local_ref" ]] && continue
            if [[ "$local_ref" =~ ^refs/tags/(.+) ]]; then
                pushed_tag="${BASH_REMATCH[1]}"
                break
            fi
        done
    fi

    if [[ -n "$pushed_tag" ]]; then
        local tag_annotation
        tag_annotation="$(git tag -l --format='%(contents)' "$pushed_tag" 2>/dev/null | head -n 1 || true)"
        hook_pre_push_tag "$pushed_tag" "$tag_annotation"
    else
        hook_pre_push_branch
    fi
}

# ---------------------------------------------------------------------------
# Tier 3a: Branch push (Build System Sync & Tests)
# Ref: docs/BEFORE_RELEASE.md (Section 2: Build System Sync, Section 3: Testing)
# Verifies Bazel/Gazelle module synchronization, unit tests, and smoke tests.
# ---------------------------------------------------------------------------
hook_pre_push_branch() {
    info "Tier 3: Running build system synchronization and test suites..."

    local push_diff
    push_diff="$(git diff --name-only origin/main...HEAD 2>/dev/null || true)"
    [[ -z "$push_diff" ]] && push_diff="$(git diff --name-only HEAD~1 2>/dev/null || true)"

    local has_go_or_build=false has_cmd_or_pkg=false has_bazel=false
    command -v bazel >/dev/null 2>&1 && has_bazel=true
    { echo "$push_diff" | grep -qE '(\.go$|\.bazel$|^go\.(mod|sum)$)' && has_go_or_build=true; } || true
    { echo "$push_diff" | grep -qE '^(cmd/|pkg/|main\.go|scripts/command-test\.sh|^go\.(mod|sum)$)' && has_cmd_or_pkg=true; } || true

    if [[ "$has_go_or_build" == "false" && "$has_cmd_or_pkg" == "false" ]]; then
        info "No Go code or build definitions modified in pushed commits; skipping tests."
        return 0
    fi

    if [[ "$has_go_or_build" == "true" ]]; then
        if [[ "$has_bazel" == "true" ]]; then
            info "Syncing Bazel and Gazelle build definitions..."
            run_quiet bazel run @rules_go//go -- mod tidy -v
            run_quiet bazel run //:gazelle
            run_quiet bazel mod tidy

            if ! git diff --exit-code -- '*BUILD.bazel' MODULE.bazel.lock >/dev/null 2>&1; then
                abort "Bazel build definitions changed after sync. Please commit generated build files before pushing."
            fi
            success "Bazel modules and Gazelle build files are in sync."

            run_step "Bazel test suite" bazel test //...
        else
            warn "'bazel' not found in PATH; falling back to 'go test ./...'."
            run_step "Go unit tests" go test ./...
        fi
    fi

    if [[ "$has_cmd_or_pkg" == "true" ]]; then
        run_step "CLI smoke tests" ./scripts/command-test.sh
    fi
}

# ---------------------------------------------------------------------------
# Tier 3b: Tag push (Pre-Release Checklist)
# Ref: docs/BEFORE_RELEASE.md (Section 4: Docs & Metadata, Section 5: Working Tree & Release Tag)
# Verifies repository state, tag format, main branch sync, and commit history.
# ---------------------------------------------------------------------------
hook_pre_push_tag() {
    local tag_name="${1:-}" tag_msg="${2:-}"

    info "Checking pre-release readiness checklist..."

    if [[ -n "$tag_name" ]]; then
        info "Validating tag name format: $tag_name"
        if ! [[ "$tag_name" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
            abort "Release tag '$tag_name' must follow SemVer and start with 'v' (e.g. v1.0.0, v0.5.2-beta1)."
        fi
        success "Tag name '$tag_name' is valid."

        if [[ -n "$tag_msg" ]]; then
            success "Tag message provided: '$tag_msg'"
        fi
    fi

    info "Checking git working directory cleanliness..."
    if [[ -n "$(git status --porcelain)" ]]; then
        abort "Working directory has unstaged or untracked changes. Ensure 'git status' is completely clean before release."
    fi
    success "Working directory is completely clean."

    info "Verifying sensitive files are not tracked in git..."
    if [[ -n "$(git ls-files "${SENSITIVE_FILES[@]}" 2>/dev/null)" ]]; then
        abort "Sensitive file(s) are tracked by Git."
    fi
    success "No sensitive files are tracked."

    if [[ "$(git branch --show-current)" != "main" ]]; then
        abort "Releases must be tagged from the 'main' branch (currently on '$(git branch --show-current)')."
    fi

    info "Verifying sync with origin/main..."
    git fetch origin main >/dev/null 2>&1 || warn "Could not fetch from origin. Check network connection."
    if git rev-parse --verify origin/main >/dev/null 2>&1; then
        local counts ahead behind
        counts="$(git rev-list --left-right --count HEAD...origin/main 2>/dev/null || echo "0 0")"
        ahead="${counts%%	*}"
        behind="${counts##*	}"
        if [[ "$behind" -gt 0 ]]; then
            abort "Local branch 'main' is behind origin/main by $behind commit(s). Please run 'git pull origin main'."
        elif [[ "$ahead" -gt 0 ]]; then
            warn "Local branch 'main' is ahead of origin/main by $ahead commit(s). Remember to push your commits."
        else
            success "Branch 'main' is synchronized with origin/main."
        fi
    fi

    info "Verifying documentation and metadata files..."
    for doc in "${DOC_FILES[@]}"; do
        if [[ ! -f "$doc" ]]; then
            abort "Missing expected documentation or metadata file: $doc"
        fi
    done
    success "Documentation and metadata files verified."

    local last_tag
    last_tag="$(git describe --tags --abbrev=0 2>/dev/null || true)"
    if [[ -n "$last_tag" ]]; then
        info "Checking commit messages since $last_tag..."
        local invalid_commits=0
        while IFS= read -r commit_line; do
            [[ -z "$commit_line" ]] && continue
            if ! is_gitmoji "${commit_line#* }"; then
                warn "Commit ${commit_line%% *} does not follow gitmoji: '${commit_line#* }'"
                invalid_commits=$((invalid_commits + 1))
            fi
        done < <(git log "${last_tag}..HEAD" --oneline)

        if [[ $invalid_commits -eq 0 ]]; then
            success "All commits since $last_tag adhere to Gitmoji conventions."
        else
            warn "$invalid_commits commit(s) do not start with a Gitmoji."
        fi
    fi

    success "Pre-release checklist verified. Ready for tagging and release!"
}

# ---------------------------------------------------------------------------
# Hook Installation Helper
# ---------------------------------------------------------------------------
install_git_hooks() {
    local hooks_dir=".git/hooks"
    if [[ ! -d "$hooks_dir" ]]; then
        abort "Directory '$hooks_dir' not found. Run this from within the git repository."
    fi

    info "Installing git hooks into $hooks_dir..."

    for hook in commit-msg pre-commit pre-push; do
        cat > "$hooks_dir/$hook" << 'EOF'
#!/usr/bin/env bash
./scripts/pre-hooks.sh "$(basename "$0")" "$@"
EOF
        chmod +x "$hooks_dir/$hook"
    done

    success "Git hooks installed: commit-msg, pre-commit, pre-push."
}

# ---------------------------------------------------------------------------
# CLI Dispatcher
# ---------------------------------------------------------------------------
usage() {
    cat << EOF
Usage: $(basename "$0") <hook_type> [arguments...]

Tiered Git pre-hooks based on docs/BEFORE_RELEASE.md.

Hook Types:
  pre-commit                Code hygiene, sensitive files, resource registration, license & lint
  commit-msg [file|text]    Validate commit message adheres to Gitmoji convention
  pre-push [vX.Y.Z] [msg]   Run build sync & tests (branch) or pre-release checklist (tag)

Commands:
  all [vX.Y.Z] [msg]        Run pre-commit, commit-msg, and pre-push in sequence
  install                   Install hooks into .git/hooks (pre-commit, commit-msg, pre-push)
  help                      Show this help message

Examples:
  # Run pre-commit checks before making a commit
  $(basename "$0") pre-commit

  # Validate commit message format
  $(basename "$0") commit-msg ":sparkles: Add video caption tool"

  # Run build sync and test suites before pushing a branch
  $(basename "$0") pre-push

  # Run pre-release checks for a tag
  $(basename "$0") pre-push v0.5.2 "🛶"
EOF
}

HOOK_TYPE="${1:-}"
shift 2>/dev/null || true

case "$HOOK_TYPE" in
    pre-commit)
        hook_pre_commit
        ;;
    commit-msg)
        hook_commit_msg "$@"
        ;;
    pre-push)
        hook_pre_push "$@"
        ;;
    all)
        hook_pre_commit
        hook_commit_msg
        hook_pre_push_branch
        hook_pre_push_tag "$@"
        ;;
    install)
        install_git_hooks
        ;;
    help|--help|-h|"")
        usage
        ;;
    *)
        error "Unknown hook type: '$HOOK_TYPE'"
        usage
        exit 1
        ;;
esac
