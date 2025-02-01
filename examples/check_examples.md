# Check Command Examples

## Outdated Modules

Check for outdated modules:
```bash
v2b check --type outdated
```

Output:
```
Outdated Modules:
+--------------------------------+---------+-------------+
|             MODULE             | CURRENT |   LATEST   |
+--------------------------------+---------+-------------+
| github.com/example/mod1        | v1.2.0  | v1.3.0     |
| github.com/example/mod2        | v2.0.0  | v2.1.0     |
+--------------------------------+---------+-------------+
```

## Security Vulnerabilities

Check for security vulnerabilities:
```bash
v2b check --type security
```

Output:
```
Security Vulnerabilities:
+--------------------------------+----------+------------------+--------------------+
|             MODULE             | SEVERITY |       CVE       |    DESCRIPTION    |
+--------------------------------+----------+------------------+--------------------+
| github.com/example/mod1        | HIGH     | CVE-2024-1234   | Buffer overflow   |
| github.com/example/mod2        | MEDIUM   | CVE-2024-5678   | SQL injection     |
+--------------------------------+----------+------------------+--------------------+
```

## Deprecated Modules

Check for deprecated modules:
```bash
v2b check --type deprecated
```

Output:
```
Deprecated Modules:
+--------------------------------+---------+----------------------------------------+
|             MODULE             | VERSION |              ALTERNATIVE               |
+--------------------------------+---------+----------------------------------------+
| github.com/example/old-mod     | v1.0.0  | Use github.com/example/new-mod instead|
+--------------------------------+---------+----------------------------------------+
```

## Version Compatibility

Check version compatibility:
```bash
v2b check --type compatibility
```

Output:
```
Incompatible Modules:
+--------------------------------+---------+------------------+
|             MODULE             | VERSION |      ISSUE      |
+--------------------------------+---------+------------------+
| github.com/example/mod1        | v1.2.0  | Requires Go 1.18|
+--------------------------------+---------+------------------+
```

## Module Cleanup

Check for unused modules:
```bash
v2b check --type cleanup
```

Output:
```
Unused Modules:
+--------------------------------+---------+
|             MODULE             | VERSION |
+--------------------------------+---------+
| github.com/example/unused-mod  | v1.0.0  |
+--------------------------------+---------+

Run with --auto-fix to remove unused modules
```

Automatically remove unused modules:
```bash
v2b check --type cleanup --auto-fix
```

## Combined Checks

Check all aspects:
```bash
v2b check --type all
```

This will run all checks in sequence and provide a comprehensive report. 