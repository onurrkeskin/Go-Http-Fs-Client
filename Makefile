VERSION := 1.0

all: fs-server fs-client

fs-server:
	docker build \
		-f Docker/dockerfile.file-server \
		-t file-server-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

fs-client:
	docker build \
		-f Docker/dockerfile.file-client \
		-t file-client-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

fs-server-run: fs-server
	docker run -p 8081:8081 file-server-amd64:$(VERSION)

fs-client-run: fs-client
	docker run -p 8080:8080 file-client-amd64:$(VERSION)