package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

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
	tnFile, tnHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get thumbnail file from form", err)
		return
	}
	tnType := tnHeader.Header.Get("Content-Type")

	// Read all image data into a byte slice
	tnData, err := io.ReadAll(tnFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't read thumbnail file", err)
		return
	}

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

	// Encode the thumbnail data to base64 to store in database
	tnString := base64.StdEncoding.EncodeToString(tnData)
	tnUrl := fmt.Sprintf("data:%s;base64,%s", tnType, tnString)

	// Update the video metadata for the new thumbnail
	video.ThumbnailURL = &tnUrl
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video with thumbnail URL", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
	// END: akhdanfadh's implementation of multipart uploads
}
