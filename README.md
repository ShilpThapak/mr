## Usage:

### Sequential:
```
go build -buildmode=plugin -race -o plugins/wc.so plugins/wc/wc.go
./bin/sequential plugins/wc/wc.so inputs/pg-*.txt
```

### Distributed:
Start Cordinator:
```
go build -buildmode=plugin -race -o plugins/wc/wc.so plugins/wc/wc.go
./bin/cordinator inputs/pg-*.txt
```
Start Workers:
```
go build -buildmode=plugin -race -o plugins/wc/wc.so plugins/wc/wc.go
./bin/worker plugins/wc/wc.so
```

---
Don't forget to clean up temp files before retrying
```
rm -r outputs/mr-*.txt intermediate/mr-*.txt
```

### View Output
```
cat outputs/mr-*.txt | sort | more
```

### Pending
- Make NReduce variable
- Add Tests to check for fault tolerance