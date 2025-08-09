package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// START: akhdanfadh's implementation of multipart form uploads

	// Parse the multipart form data into the request object itself
	// with large enough memory to hold in RAM, rest will be stored in a temporary file
	const maxMemory = 10 << 20 // left bit shift; 10 in binary is 1010, so 1010 with 20 trailing zeros (10MB)
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse multipart form", err)
		return
	}

	// Get image data from the form
	tnData, tnHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get thumbnail file from form", err)
		return
	}
	tnType := tnHeader.Header.Get("Content-Type")

	// Get the video metadata from the database
	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get video", err)
		return
	}
	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "You are not allowed to upload thumbnail for this video", nil)
		return
	}

	// Save the thumbnail to a file
	tnFileExt, ok := strings.CutPrefix(tnType, "image/")
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid thumbnail content type", nil)
		return
	}
	tnFileName := fmt.Sprintf("%s.%s", videoID, tnFileExt)
	tnFilePath := filepath.Join(cfg.assetsRoot, tnFileName)
	tnFile, err := os.Create(tnFilePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create thumbnail file", err)
		return
	}
	_, err = io.Copy(tnFile, tnData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't write thumbnail file", err)
		return
	}

	// Update the video metadata for the new thumbnail
	tnUrl := fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, tnFileName)
	video.ThumbnailURL = &tnUrl
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video with thumbnail URL", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
	// END: akhdanfadh's implementation of multipart uploads
}
