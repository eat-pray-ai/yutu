# How to Test

This guide outlines the testing strategy for `yutu`. We combine high-level strategic goals with specific patterns to ensure code reliability and maintainability.

## 0. High-Level Strategy

### 0.1. Philosophy
Our testing philosophy focuses on **confidence** and **maintainability**. Tests should give us confidence that the code works as expected and allow us to refactor safely. We prefer comprehensive unit tests that mock external dependencies (YouTube API) over flaky end-to-end tests.

### 0.2. Test-Driven Development (TDD)
We encourage a TDD approach, especially when adding new resources:
1.  **Define the Interface**: Write the `TestNew<Resource>` test first. This forces you to design the API surface (Options pattern) before implementation.
2.  **Define Behavior**: Write the test for a specific method (e.g., `Get` or `Insert`) using `httptest` to mock the expected API interaction.
3.  **Implement**: Write the code to make the tests pass.
4.  **Refactor**: Clean up the code while keeping tests green.

### 0.3. Coverage Targets
- **Target**: We aim for **>80%** code coverage for domain logic in `pkg/`.
- **Critical Paths**: Authentication, flag parsing, and API request construction must have near 100% coverage.
- **Tools**: Use `go test -cover` or the Bazel coverage report to identify gaps.

### 0.4. Continuous Integration
Tests are run automatically on every Pull Request via GitHub Actions.
- **Fast Feedback**: Unit tests (`go test ./...`) should be fast.
- **Hermeticity**: Tests should not depend on external internet access (hence `httptest` for API calls).

---

## 1. Test Location & Naming

- **Location**: Tests must be co-located with the code in `pkg/<resource>/<resource>_test.go`.
- **Naming**: Use `Test<Resource>_<Method>` (e.g., `TestChannel_Get`, `TestPlaylist_Insert`).
- **Constructor**: `TestNew<Resource>` tests the option pattern implementation.

## 2. Testing Constructors (`New<Resource>`)

We use the Functional Options pattern. Tests must ensure that all options are applied correctly and edge cases are handled.

**Required Test Cases:**
- `with all options`: Verify every field is set correctly.
- `with no options`: Verify defaults.
- `with nil/false boolean options`: Ensure nil pointers and false values are handled distinctively.
- `with zero/negative max results`: Verify boundary logic (e.g., `0` -> `math.MaxInt64`, `<0` -> `1`).
- `with empty string values`: Ensure empty strings don't crash or set incorrect defaults if not intended.

**Example Structure:**
```go
func TestNewChannel(t *testing.T) {
    type args struct {
        opts []Option
    }
    tests := []struct {
        name string
        args args
        want IChannel[youtube.Channel]
    }{
        {
            name: "with all options",
            args: args{opts: []Option{WithTitle("Test"), ...}},
            want: &Channel{Title: "Test", ...},
        },
        // ... other cases
    }
    // Run loop using reflect.DeepEqual
}
```

## 3. Testing Methods (`Get`, `List`, `Insert`, `Update`, `Delete`)

We **do not** use a mocking library for the YouTube service. Instead, we use `httptest.NewServer` to mock the Google API endpoint. This allows us to verify the exact HTTP requests (method, query params, body) that the library sends.

### 3.1. Mocking the API

Create a test server that acts as the YouTube API.

```go
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // 1. Verify Request
    if r.Method != "GET" {
        t.Errorf("expected GET, got %s", r.Method)
    }
    if r.URL.Query().Get("part") != "snippet" {
        t.Errorf("expected part=snippet")
    }

    // 2. Mock Response
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{ "items": [ { "id": "1", "snippet": { "title": "Test" } } ] }`))
}))
defer ts.Close()
```

### 3.2. Injecting the Mock Service

Initialize the `youtube.Service` with the mock server's URL.

```go
svc, _ := youtube.NewService(
    context.Background(),
    option.WithEndpoint(ts.URL), // Point to local test server
    option.WithAPIKey("test-key"),
)
```

### 3.3. Pagination

For `Get` methods that support pagination, simulate `pageToken` logic in the mock handler.

```go
func(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("pageToken")
    if token == "" {
        // Return page 1 and nextPageToken
    } else if token == "page-2" {
        // Return page 2
    }
}
```

## 4. Testing Output Formats

For `List` methods, verify that all supported output formats (JSON, YAML, Table) work without error and write to the buffer.

```go
tests := []struct {
    name   string
    output string // "json", "yaml", "table"
}{
    { "list json", "json" },
    { "list yaml", "yaml" },
    { "list table", "table" },
}
```

## 5. Running Tests

You can run tests using standard Go tools or Bazel.

**Standard Go:**
```bash
go test ./pkg/...
# Verbose single test
go test -v ./pkg/channel -run TestChannel_Get
```

**Bazel:**
```bash
bazel test //pkg/...
```
