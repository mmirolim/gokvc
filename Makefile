APP = kvc
# set commit short
COMMIT =`git rev-parse --short HEAD`
# set build time
TIMESTM = `date -u '+%Y-%m-%d_%H:%M:%S%p'`
# format version signature
FORMAT = $(COMMIT)-$(TIMESTM)

build:
	go build -o $(APP) -ldflags "-X main.BuildVersion=$(FORMAT)"
run: build
	./$(APP) -log_dir="logs" -stderrthreshold=INFO -v=5
