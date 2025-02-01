# Search Command Examples

## Content-Based Search

Search for modules containing "http" in their content:
```bash
v2b search --content "http"
```

Output:
```
Found 3 modules matching "http":

1. github.com/example/http-client (v1.2.0)
   Match: pkg/client/http.go:15 - "HTTP client implementation"

2. github.com/example/server (v2.0.0)
   Match: internal/server/http_handler.go:25 - "HTTP request handler"

3. github.com/example/utils (v0.5.0)
   Match: utils/http_helpers.go:10 - "HTTP utility functions"
```

## Version-Based Search

Search for modules with specific version:
```bash
v2b search --version "v1.0.0"
```

Search for modules with version range:
```bash
v2b search --version ">=v1.0.0 <v2.0.0"
```

## Date-Based Search

Search modules updated after a specific date:
```bash
v2b search --from "2024-01-01"
```

Search modules within a date range:
```bash
v2b search --from "2024-01-01" --to "2024-02-01"
```

## Combined Search

Search for HTTP-related modules with specific version:
```bash
v2b search --content "http" --version "v1.*"
```

Search with date range and content:
```bash
v2b search --content "database" --from "2024-01-01" --to "2024-02-01"
```

## Pagination

Get first page of results:
```bash
v2b search --content "http" --page 1 --limit 5
```

Get next page:
```bash
v2b search --content "http" --page 2 --limit 5
```

## Search Result Highlighting

Search with highlighted matches:
```bash
v2b search --content "server" --highlight
```

Output:
```
Found match in github.com/example/server:
Line 15: "A high-performance [server] implementation"
Line 25: "Configure [server] options and middleware"
``` 