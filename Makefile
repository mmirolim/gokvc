APP = gokvc
# set commit short
COMMIT =`git rev-parse --short HEAD`
# set build time
TIMESTM = `date -u '+%Y-%m-%d_%H:%M:%S%p'`
# format version signature
FORMAT = $(COMMIT)-$(TIMESTM)

bench-cache:
	go test ./cache -run=none -bench=^Benchmark -cpu=1,2,3 -cpuprofile=prof.cpu -memprofile=prof.mem

bench-api:
	go test ./api -run=none -bench=^Benchmark -cpu=1,2,3 -cpuprofile=prof.cpu -memprofile=prof.mem

test:
	go test -v ./...

test-race:
	go test -v -race ./...

contention-prof:
	go test ./... -bench=Parallel -blockprofile=prof.block

get-deps:
	go get ./...

build: get-deps
	mkdir -p logs
	go build -o $(APP) -i -ldflags "-X main.BuildVersion=$(FORMAT)"

build-race:  get-deps
	mkdir -p logs
	go build -o $(APP) -race -ldflags "-X main.BuildVersion=$(FORMAT)"

run: build
	./$(APP) -log_dir="logs" -stderrthreshold=INFO -v=5

run-race: build-race
	./$(APP) -log_dir="logs" -stderrthreshold=INFO -v=5
