# ptb

BINARY   := ptb
SRC      := main.go
TXT_DIR  := txt
OUT_DIR  := output
PORT     := 8080

.PHONY: all build generate serve clean

all: build generate

build:
	go build -o $(BINARY) $(SRC)

generate: build
	@if [ ! -f $(OUT_DIR)/index.html ] || \
	   [ -n "$$(find $(TXT_DIR) -type f -newer $(OUT_DIR)/index.html)" ]; then \
	  echo "[INFO] Changes detected --> regenerating site…"; \
	  ./$(BINARY); \
	else \
	  echo "[INFO] No changes in $(TXT_DIR), skipping generation."; \
	fi

serve: generate
	@echo "Serving at http://localhost:$(PORT)/"
	cd $(OUT_DIR) && python3 -m http.server $(PORT)

clean:
	rm -f $(BINARY)

