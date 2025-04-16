.PHONY: all

ENTRY_POINT := "./cmd/entrypoint"
PROGRAM := "./authService"
LOCAL := "--config=./config/local.yaml"
PROD := "--config=./config/prod.yaml"
FILES_DIR := "./files"

all: build run_local

build:
	@echo Compile and build...
	@go build -o $(PROGRAM) $(ENTRY_POINT)

run_local:
	@echo Run local app: $(PROGRAM) $(LOCAL)
	@$(PROGRAM) $(LOCAL)

run_prod:
	@echo Run prod app: $(PROGRAM) $(PROD)
	@$(PROGRAM) $(PROD)

clear:
	@echo Cleaning files...
	@rm -rf $(FILES_DIR) $(EXE_NAME)