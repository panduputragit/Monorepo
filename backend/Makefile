SERVICES := api-gateway auth-service employee-service branch-service member-service membership-service attendance-service notification-service payment-service

BUILD_TARGETS := $(addprefix build-,$(SERVICES))
VET_TARGETS := $(addprefix vet-,$(SERVICES))
RUN_TARGETS := $(addprefix run-,$(SERVICES))

.PHONY: gen build-all vet-all run-all deps help \
	$(BUILD_TARGETS) $(VET_TARGETS) $(RUN_TARGETS)

# Generate proto code
gen:
	@echo "Generating proto code..."
	cd protos && buf generate

# -----------------------
# Build
# -----------------------
build-all: $(BUILD_TARGETS)

$(BUILD_TARGETS):
	@service=$(@:build-%=%); \
	echo "Building $$service..."; \
	cd app/$$service && go build -buildvcs=false -o ../../bin/$$service ./cmd/server

# -----------------------
# Vet
# -----------------------
vet-all: $(VET_TARGETS)

$(VET_TARGETS):
	@service=$(@:vet-%=%); \
	echo "Vetting $$service..."; \
	cd app/$$service && go vet ./...

# -----------------------
# Run
# -----------------------
run-all:
	@mkdir -p logs
	@for service in $(SERVICES); do \
		echo "Starting $$service..."; \
		(cd app/$$service && go run ./cmd/server > ../../logs/$$service.log 2>&1 &) ; \
	done
	@echo "All services started."

$(RUN_TARGETS):
	@service=$(@:run-%=%); \
	echo "Running $$service..."; \
	cd app/$$service && go run ./cmd/server

# -----------------------
# Dependencies
# -----------------------
deps:
	@for service in packages $(SERVICES:%=app/%); do \
		echo "Tidying $$service..."; \
		(cd $$service && go mod tidy); \
	done

stop-all:
	@for pidfile in pids/*.pid; do \
		if [ -f "$$pidfile" ]; then \
			pid=$$(cat $$pidfile); \
			echo "Stopping PID $$pid..."; \
			kill $$pid 2>/dev/null || true; \
			rm -f $$pidfile; \
		fi; \
	done
	@echo "All services stopped."

help:
	@echo "Available targets:"
	@echo "  gen          - Generate proto code"
	@echo "  build-all    - Build all services"
	@echo "  build-<name> - Build one service"
	@echo "  vet-all      - Vet all services"
	@echo "  vet-<name>   - Vet one service"
	@echo "  run-all      - Run all services"
	@echo "  run-<name>   - Run one service"
	@echo "  deps         - Run go mod tidy"