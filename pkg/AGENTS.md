# Package Guidelines

## OVERVIEW

Core domain logic and shared infrastructure.

## STRUCTURE

- `<resource>/`: Domain logic (video, channel, etc.).
- `auth/`: Authentication service.
- `utils/`: Generic helpers.
- `common/`: Shared types.

## PATTERNS

- 1:1 mapping with YouTube API resources.
- Structure: `<resource>.go` + `<resource>_test.go` + `BUILD.bazel`.
- Domain packages depend on `auth` and `common`.

## TESTING

- Table-driven tests.
- Real `youtube.Service` with `httptest.NewServer`.
- Global state (envs) carefully managed/restored.
