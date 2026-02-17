# ptb

BINARY   := ptb
PKG      := .
TXT_DIR  := txt
TPL_DIR  := templates
OUT_DIR  := output
PORT     := 8080

.PHONY: all build generate generate-ci serve clean

all: build generate

build:
	go build -o $(BINARY) $(PKG)

generate: build
	@if [ ! -f $(OUT_DIR)/index.html ] || \
	   [ ! -f $(OUT_DIR)/rss.xml ] || \
	   [ ./$(BINARY) -nt $(OUT_DIR)/index.html ] || \
	   [ -n "$$(find $(TXT_DIR) -type f -newer $(OUT_DIR)/index.html)" ]; then \
	  echo "[INFO] Build/input changes detected --> regenerating site…"; \
	  ./$(BINARY); \
	elif [ -n "$$(find $(TPL_DIR) -type f -newer $(OUT_DIR)/index.html)" ]; then \
	  echo "[INFO] Changes detected --> regenerating site…"; \
	  ./$(BINARY); \
	else \
	  echo "[INFO] No changes in $(TXT_DIR) or $(TPL_DIR), skipping generation."; \
	fi

generate-ci: build
	@echo "[INFO] CI mode: regenerating site deterministically..."
	./$(BINARY)

serve: generate
	@echo "Serving at http://localhost:$(PORT)/"
	cd $(OUT_DIR) && python3 -m http.server $(PORT)

clean:
	rm -f $(BINARY)
