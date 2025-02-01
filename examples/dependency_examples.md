# Dependency Command Examples

## Pin Module

Pin a module to its current version:
```bash
v2b dependency pin --module github.com/example/mod
```

Output:
```
Successfully pinned github.com/example/mod to version v1.2.0
```

## Unpin Module

Remove version pin from a module:
```bash
v2b dependency unpin --module github.com/example/mod
```

Output:
```
Successfully unpinned github.com/example/mod
```

## Explain Dependencies

Show why a module is needed:
```bash
v2b dependency why --module github.com/example/mod
```

Output:
```
Dependencies for github.com/example/mod:
+--------------------------------+---------+--------------------------------+
|             MODULE             | VERSION |             REASON             |
+--------------------------------+---------+--------------------------------+
| github.com/example/dep1        | v1.0.0  | Required by internal/server   |
| github.com/example/dep2        | v2.0.0  | Required by pkg/client        |
+--------------------------------+---------+--------------------------------+
```

## Dependency Graph

Show dependency graph for a module:
```bash
v2b dependency graph --module github.com/example/mod
```

Output:
```
Dependency graph for github.com/example/mod:

github.com/example/mod@v1.2.0
  └── github.com/example/dep1@v1.0.0
      └── github.com/example/subdep1@v0.5.0
      └── github.com/example/subdep2@v1.1.0
  └── github.com/example/dep2@v2.0.0
      └── github.com/example/subdep3@v0.8.0
```

## Dependency Analysis

Analyze dependencies for potential issues:
```bash
v2b dependency analyze --module github.com/example/mod
```

Output:
```
Dependency Analysis:

Direct Dependencies: 2
Indirect Dependencies: 3
Total Size: 5.2 MB

Issues Found:
- Circular dependency between dep1 and dep2
- Multiple versions of subdep1 in use
- Unused dependency: subdep3

Recommendations:
1. Resolve circular dependency
2. Consolidate subdep1 versions
3. Consider removing unused dependency
```

## Dependency Updates

Check for available dependency updates:
```bash
v2b dependency updates --module github.com/example/mod
```

Output:
```
Available Updates:
+--------------------------------+---------+---------+-------------+
|             MODULE             | CURRENT | LATEST  |    TYPE     |
+--------------------------------+---------+---------+-------------+
| github.com/example/dep1        | v1.0.0  | v1.1.0  | Minor       |
| github.com/example/subdep2     | v1.1.0  | v2.0.0  | Major       |
+--------------------------------+---------+---------+-------------+
```

## Dependency Tree

Show full dependency tree with details:
```bash
v2b dependency tree --module github.com/example/mod
```

Output:
```
github.com/example/mod@v1.2.0
├── github.com/example/dep1@v1.0.0
│   ├── github.com/example/subdep1@v0.5.0
│   │   └── [license: MIT, size: 1.2MB]
│   └── github.com/example/subdep2@v1.1.0
│       └── [license: Apache-2.0, size: 850KB]
└── github.com/example/dep2@v2.0.0
    └── github.com/example/subdep3@v0.8.0
        └── [license: BSD-3, size: 2.1MB]
``` 