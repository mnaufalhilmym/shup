package controller

import (
	"net/http"
	"path"
	"strconv"

	"github.com/google/uuid"
	"shup.hilmy.dev/src/internal/service"
)

type File struct {
	svc *service.File
}

func NewFile(svc *service.File) *File {
	return &File{svc}
}

func (c *File) Save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileName := path.Base(r.URL.Path)
	if fileName == "" {
		fileName = uuid.NewString()
	}

	file, err := c.svc.Save(ctx, fileName, r.Body)
	if err != nil {
		http.Error(w, "Failed to save the file", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Uploaded 1 file, " + strconv.FormatInt(file.Size, 10) + " bytes\n\n" +
		"wget https://shup.hilmy.dev/" + file.ID.String()))
}

func (c *File) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileIDStr := path.Base(r.URL.Path)
	if fileIDStr == "" {
		http.Error(w, "FileID not provided", http.StatusBadRequest)
		return
	}

	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		http.Error(w, "FileID not valid", http.StatusBadRequest)
		return
	}

	file, err := c.svc.Get(ctx, fileID)
	if err != nil {
		http.Error(w, "Failed to get the file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.FileName+"\"")

	// Serve the file content
	http.ServeFile(w, r, "./files/"+file.ID.String())
}
