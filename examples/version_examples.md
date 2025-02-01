# Version Command Examples

## Version Upgrade

Upgrade a module to a specific version:
```bash
v2b version upgrade --module github.com/example/mod --version v1.2.0
```

Output:
```
Successfully upgraded github.com/example/mod to version v1.2.0
```

## Version Rollback

Rollback a module to its previous version:
```bash
v2b version rollback --module github.com/example/mod
```

Output:
```
Successfully rolled back github.com/example/mod from v1.2.0 to v1.1.0
```

## Version Conflict Resolution

Resolve version conflicts automatically:
```bash
v2b version resolve
```

Output:
```
Analyzing version conflicts...

Found conflicts:
1. github.com/example/mod1
   - Required v1.0.0 by module A
   - Required v1.1.0 by module B
   Resolution: Upgraded to v1.1.0

2. github.com/example/mod2
   - Required v2.0.0 by module C
   - Required v2.1.0 by module D
   Resolution: Upgraded to v2.1.0

All conflicts resolved successfully
```

## Version History

View version history for a module:
```bash
v2b version history --module github.com/example/mod
```

Output:
```
Version History for github.com/example/mod:
+----------+------------+------------------+
| VERSION  |    DATE    |     CHANGES     |
+----------+------------+------------------+
| v1.2.0   | 2024-02-20 | Feature release |
| v1.1.0   | 2024-02-10 | Bug fixes       |
| v1.0.0   | 2024-02-01 | Initial release |
+----------+------------+------------------+
```

## Version Comparison

Compare two versions of a module:
```bash
v2b version compare --module github.com/example/mod --from v1.0.0 --to v1.1.0
```

Output:
```
Changes between v1.0.0 and v1.1.0:

Added:
- New feature X
- Support for Y

Fixed:
- Bug in function Z
- Performance issue in A

Breaking Changes:
- Renamed function B to C
```

## Version Tracking

Track version changes over time:
```bash
v2b version track --module github.com/example/mod
```

Output:
```
Version Changes (Last 30 days):
+------------+----------+----------+------------------+
|    DATE    |   FROM   |    TO    |     REASON      |
+------------+----------+----------+------------------+
| 2024-02-20 | v1.1.0   | v1.2.0   | Manual upgrade  |
| 2024-02-15 | v1.2.0   | v1.1.0   | Rollback        |
| 2024-02-10 | v1.0.0   | v1.1.0   | Auto-upgrade    |
+------------+----------+----------+------------------+
``` 