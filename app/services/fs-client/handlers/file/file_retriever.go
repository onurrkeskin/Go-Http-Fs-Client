package file

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/onurkeskin/go-http-fs-client/domain/core/file_service"
	"gitlab.com/onurkeskin/go-http-fs-client/foundation/web"
)

type FileRetriverHandler struct {
	fileService file_service.FileService
}

func (fileRetrieverHandller FileRetriverHandler) RetrieveFiles(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	wantedString := web.Param(r, "token")

	files, err := fileRetrieverHandller.fileService.GetFiles(ctx, wantedString)
	if err != nil {
		return fmt.Errorf("Unable to download files: %w", err)
	}

	return web.Respond(ctx, w, files, http.StatusOK)
}
