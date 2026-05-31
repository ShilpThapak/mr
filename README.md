## Usage:

### Sequential:
```
go run cmd/sequential/main.go
```

### Distributed:
Start Cordinator:
```
go run cmd/cordinator/main.go
```
Start Workers:
```
go run cmd/worker/main.go
```

---
Don't forget to clean up temp files before restating
```
rm -r outputs/mr-*.txt intermediate/mr-*.txt
```

### View Output
```
cat outputs/mr-*.txt | sort | more
```

### Pending
- Add Gorotines to make it parallel
- Add support for custom map reduce functions
- Make NReduce variable
- Improve the cli calling lines
- All Shell script to run this multiple times
- Add Tests to check for fault tolerance