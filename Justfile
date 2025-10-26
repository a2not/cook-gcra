bench:
  go test -bench=. -benchmem ./... -benchtime=10s -cpu=1,16,256
