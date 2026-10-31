.PHONY: help check api-check api-build admin-build web-build build fmt

help:
	@echo "Targets:"
	@echo "  make check        - API vet+test, admin+web production build"
	@echo "  make api-check    - go vet + go test (qixue-api)"
	@echo "  make api-build    - build qixue-api → qixue-api/bin/server"
	@echo "  make admin-build  - npm ci + build (qixue-admin)"
	@echo "  make web-build    - npm ci + build (qixue-web)"
	@echo "  make build        - api-build + admin-build + web-build"
	@echo "  make fmt          - gofmt API sources"

check: api-check admin-build web-build

api-check:
	$(MAKE) -C qixue-api check

api-build:
	$(MAKE) -C qixue-api build

fmt:
	$(MAKE) -C qixue-api fmt

admin-build:
	cd qixue-admin && npm ci --ignore-scripts && npm run build

web-build:
	cd qixue-web && npm ci --ignore-scripts && npm run build

build: api-build admin-build web-build
