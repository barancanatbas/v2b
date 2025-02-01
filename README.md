# v2b - Go Module Manager

v2b is a powerful Go module management tool that helps you manage your project's dependencies with advanced features and a user-friendly interface.

## Features

- 📋 **Enhanced Listing**: View detailed module information including size, usage, license, and last update
- 🔍 **Advanced Search**: Search modules by content, version, or date range
- 🛠️ **Module Maintenance**: Check for outdated modules, security vulnerabilities, and compatibility
- 📦 **Version Management**: Upgrade, rollback, and resolve version conflicts
- 🔗 **Dependency Management**: Pin versions, track dependencies, and visualize dependency graphs

## Installation

```bash
go install github.com/barancanatbas/v2b@latest
```

## Quick Start

1. List all modules:
```bash
v2b list
```

2. Search for modules:
```bash
v2b search --content "http" --version "v1.0.0"
```

3. Check module status:
```bash
v2b check --type outdated
v2b check --type security
```

4. Manage versions:
```bash
v2b version upgrade --module github.com/example/mod --version v1.2.0
v2b version rollback --module github.com/example/mod
```

5. Handle dependencies:
```bash
v2b dependency pin --module github.com/example/mod
v2b dependency why --module github.com/example/mod
v2b dependency graph --module github.com/example/mod
```

## Command Reference

### List Command
```bash
v2b list [flags]
  --sort string    Sort by: date, size, name (default "name")
  --filter string  Filter modules by pattern
  --format string  Output format: table, json, csv (default "table")
```

### Search Command
```bash
v2b search [flags]
  --content string   Search in module content
  --version string   Search by version
  --from string      Search from date (YYYY-MM-DD)
  --to string        Search to date (YYYY-MM-DD)
  --page int         Page number for results
  --limit int        Results per page (default 10)
```

### Check Command
```bash
v2b check [flags]
  --type string    Check type: outdated, security, deprecated, compatibility, cleanup
  --auto-fix       Automatically fix issues (only for cleanup)
```

### Version Command
```bash
v2b version [command] [flags]
Commands:
  upgrade         Upgrade module to specific version
  rollback        Rollback to previous version
  resolve         Resolve version conflicts

Flags:
  --module string   Module path to operate on
  --version string  Target version (for upgrade)
```

### Dependency Command
```bash
v2b dependency [command] [flags]
Commands:
  pin            Pin module to current version
  unpin          Remove version pin
  why            Explain module dependencies
  graph          Show dependency graph

Flags:
  --module string  Module path to operate on
```

## Configuration

v2b can be configured using a config file (`~/.v2b/config.yaml`) or environment variables:

```yaml
# ~/.v2b/config.yaml
default_format: "table"
page_size: 20
cache_dir: "~/.v2b/cache"
```

Environment variables:
- `V2B_DEFAULT_FORMAT`: Default output format
- `V2B_PAGE_SIZE`: Default page size for paginated results
- `V2B_CACHE_DIR`: Cache directory location

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
