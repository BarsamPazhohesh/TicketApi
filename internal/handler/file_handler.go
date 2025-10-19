package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/services/storage"
	"ticket-api/internal/util"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	storage *storage.StorageService
}

func NewFileHandler(storage *storage.StorageService) *FileHandler {
	return &FileHandler{
		storage: storage,
	}
}

// UploadTicketFile godoc
// @Summary      Upload a ticket file
// @Description  Upload a file to the temporary storage for a ticket. File must be multipart/form-data with field name `file`.
// @Tags         TicketFile
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Ticket file to upload"
// @Success      200   {object}  dto.IDResponse[string]  "Returns uploaded file ID"
// @Failure      400   {object}  errx.APIError
// @Failure      413   {object}  errx.APIError  "File too large"
// @Failure      415   {object}  errx.APIError  "Unsupported file extension"
// @Failure      500   {object}  errx.APIError
// @Router       /files/UploadTicketFile/ [post]
func (h *FileHandler) UploadTicketFileHandler(c *gin.Context) {
	maxUploadSize := config.Get().TicketConfig.MaxTicketUploadFileSize << 10
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			apiErr := errx.Respond(errx.ErrRequestBodyTooLarge, err)
			c.JSON(apiErr.HTTPStatus, apiErr)
			return
		}
		apiErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	file, header, err := c.Request.FormFile("file")

	if err != nil {
		apiErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	defer file.Close()

	ext, apiErr := parseTicketFileExtension(header)
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	filename := fmt.Sprintf("%s%s", util.GenerateUUID(), ext)

	_, apiErr = h.storage.UploadTicketFileToTemp(c, filename, file, header.Size, c.ContentType())
	if apiErr != nil {
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, &dto.IDResponse[string]{ID: filename})
}

// DownloadTicketFileHandler godoc
// @Summary      Download a ticket file
// @Description  Generates a temporary presigned URL for the given file and redirects the client to that URL to start the download.
// @Tags         TicketFile
// @Produce      json
// @Param        objectName  path  string  true  "File object name (UUID + extension)"
// @Param 			 ticketId body dto.IDRequest[string] true "ticket ID for file location"
// @Success      200         "Give MinIO download URL"
// @Failure      400         {object}  errx.APIError
// @Failure      404         {object}  errx.APIError  "File not found"
// @Failure      500         {object}  errx.APIError
// @Router       /files/GetDownloadLinkTicketFile/{objectName}/ [post]
func (h *FileHandler) GetDownloadLinkTicketFileHandler(c *gin.Context) {
	objectName, err := util.ParseObjectName(c.Param("objectName"))
	var req dto.IDRequest[string]
	if !bindJSON(c, &req) {
		return
	}

	if err != nil {
		apiErr := errx.Respond(errx.ErrBadRequest, err)
		c.JSON(apiErr.HTTPStatus, apiErr)
		return
	}

	url, appErr := h.storage.GetPresignedTicketFileURL(c.Request.Context(), req.ID, objectName)
	if appErr != nil {
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	c.JSON(http.StatusOK, &dto.TicketDownloadLink{Url: url})
}

// parseTicketFileExtension checks if the file's extension is supported.
// If supported, it returns the *normalized* (lowercase) extension string.
// If not supported, it returns an empty string and a specific API error.
func parseTicketFileExtension(file *multipart.FileHeader) (string, *errx.APIError) {
	allowedExts := config.Get().TicketConfig.AcceptableFilesForUpload
	ext := strings.ToLower(filepath.Ext(file.Filename))

	for _, allowed := range allowedExts {
		if strings.EqualFold(allowed, ext) {
			return ext, nil
		}
	}

	errDetail := fmt.Sprintf("file extension %q is not supported", ext)
	return "", errx.Respond(errx.ErrUnsupportedFileExtension, errors.New(errDetail))
}
