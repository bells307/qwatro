# qwatro

qwatro is a simple network tool which can scan tcp ports

# Usage
Use `--help` to see all application arguments.

To run tcp port scanning use the command:
```
qwatro scan <HOST> [flags]
```

`qwatro scan` flags:
```
-h, --help                   help for scan
-p, --port-range string      Port range (default "1-65535")
--tcp-timeout duration       Maximum response time for tcp scanning (default 300ms)
-t, --timeout duration       General scan timeout
-w, --workers int            Workers for scanning (default 200)
```

## Example
Run port scannig on range 7000-10000:
```
qwatro scan 127.0.0.1 -p 7000-10000
```

### Output
```
PS C:\dev\go\qwatro> qwatro scan 127.0.0.1 -p 7000-10000
127.0.0.1:7680
127.0.0.1:8082
```