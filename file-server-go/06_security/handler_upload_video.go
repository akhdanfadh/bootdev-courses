package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	// Parse video ID from the URL path
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid video ID", err)
		return
	}

	// Authenticate the user with bearer token
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

	// Authorize the user to upload the video
	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get video from database", err)
		return
	}
	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "You are not allowed to upload video for this video ID", nil)
		return
	}

	fmt.Println("uploading video for video ID", videoID, "by user", userID)

	// Parse the multipart form data
	const maxMemory = 1 << 30 // 1GB
	if err = r.ParseMultipartForm(maxMemory); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse multipart form", err)
		return
	}
	vidData, vidHeader, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get video file from form", err)
		return
	}
	defer vidData.Close()

	// Validate uploaded file to be mp4 format
	vidType := vidHeader.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(vidType)
	if err != nil || mediaType != "video/mp4" {
		respondWithError(w, http.StatusBadRequest, "Invalid video file type", fmt.Errorf("expected video/mp4, got %s", mediaType))
		return
	}

	// Save the video file as temporary
	tempFile, err := os.CreateTemp("", "tubely-upload.mp4")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create temporary file", err)
		return
	}
	defer os.Remove(tempFile.Name()) // Clean up temp file after upload
	defer tempFile.Close()           // note: defer is LIFO, so this order correct
	if _, err = io.Copy(tempFile, vidData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save video file", err)
		return
	}

	// Process the video file to ensure it is suitable for fast start playback
	fastPath, err := processVideoForFastStart(tempFile.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't process video for fast start", err)
		return
	}
	fastFile, err := os.Open(fastPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't open processed video file", err)
		return
	}
	defer os.Remove(fastPath) // Clean up processed file after upload
	defer fastFile.Close()

	// Prepare to upload the video file to S3
	// - Get video aspect ratio for S3 organization
	aspectRatio, err := getVideoAspectRatio(fastFile.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get video aspect ratio", err)
		return
	}
	switch aspectRatio {
	case "16:9":
		aspectRatio = "landscape"
	case "9:16":
		aspectRatio = "portrait"
	}
	// - Create random name (key) for the video file
	base := make([]byte, 32)
	if _, err = rand.Read(base); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate random file name", err)
		return
	}
	objectKey := fmt.Sprintf("%s/%s.mp4", aspectRatio, base64.RawURLEncoding.EncodeToString(base))

	// Upload the video file to S3
	objectData := &s3.PutObjectInput{
		Bucket:      aws.String(cfg.s3Bucket),
		Key:         aws.String(objectKey),
		Body:        fastFile,
		ContentType: aws.String(mediaType),
	}
	if _, err = cfg.s3Client.PutObject(r.Context(), objectData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't upload video to S3", err)
		return
	}

	// Update the video metadata for the video in database
	vidUrl := fmt.Sprintf("%s,%s", cfg.s3Bucket, objectKey) // a format we use to store unsigned video URLs
	video.VideoURL = &vidUrl
	if err = cfg.db.UpdateVideo(video); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video metadata in database", err)
		return
	}

	// since we use private bucket
	video, err = cfg.dbVideoToSignedVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get signed video URL", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}

// dbVideoToSignedVideo converts a "unsigned" video in database to a signed video with presigned URL
func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil || *video.VideoURL == "" {
		return video, nil // not error since this is a valid case for draft video
	}
	bucket, key, found := strings.Cut(*video.VideoURL, ",")
	if !found {
		return video, fmt.Errorf("video URL is not in expected unsigned format: %s", *video.VideoURL)
	}
	const expireTime = 15 * time.Minute // magic number, adjusted for bootdev video samples
	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, expireTime)
	if err != nil {
		return video, err
	}
	video.VideoURL = &presignedURL
	return video, nil
}

// generatePresignedURL generates a presigned URL for S3
func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignedClient := s3.NewPresignClient(s3Client)
	presignedHTTPRequest, err := presignedClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", fmt.Errorf("couldn't generate presigned URL: %w", err)
	}
	return presignedHTTPRequest.URL, nil
}
