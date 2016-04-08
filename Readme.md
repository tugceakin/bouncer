## Running

```
go build
./bouncer-node
```

In another terminal
```
go run test/testserver.go 9091
```

Then test using 
```
ab -n10000 -c100 localhost:9090/sdf
```
