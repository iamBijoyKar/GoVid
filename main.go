package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.StaticFile("/", "./static/index.html")
	r.Static("/static", "./static")
	r.GET("/video", streamVideo)
	r.POST("/upload", uploadVideo)
	r.GET("/videos", listVideos)
	fmt.Println("Server started at http://localhost:8080")
	fmt.Println("Upload videos at http://localhost:8080")
	r.Run(":8080")
}

func streamVideo(c *gin.Context) {
	// Get filename from query parameter, default to sample.mp4
	filename := c.Query("file")
	if filename == "" {
		filename = "sample.mp4"
	}

	// Sanitize filename to prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.String(http.StatusBadRequest, "Invalid filename")
		return
	}

	filePath := filepath.Join("videos", filename)
	file, err := os.Open(filePath)
	if err != nil {
		c.String(http.StatusNotFound, "File not found")
		return
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not get file info")
		return
	}
	fileSize := fi.Size()

	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		c.String(http.StatusBadRequest, "Range header required")
		return
	}

	parts := strings.Split(rangeHeader, "=")
	if len(parts) != 2 || parts[0] != "bytes" {
		c.String(http.StatusBadRequest, "Invalid Range header")
		return
	}

	rangeParts := strings.Split(parts[1], "-")
	start, err := strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid Range start")
		return
	}

	end := fileSize - 1
	if len(rangeParts) == 2 && rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid Range end")
			return
		}
	}

	chunkSize := end - start + 1
	buffer := make([]byte, chunkSize)
	_, err = file.ReadAt(buffer, start)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file")
		return
	}

	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", fmt.Sprintf("%d", chunkSize))
	c.Header("Content-Type", "video/mp4")
	c.Status(http.StatusPartialContent)
	c.Writer.Write(buffer)
}

func uploadVideo(c *gin.Context) {
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
		return
	}
	defer file.Close()

	// Validate file size (100MB limit)
	const maxFileSize = 100 * 1024 * 1024 // 100MB
	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 100MB"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".mp4", ".avi", ".mov", ".mkv", ".webm"}
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: mp4, avi, mov, mkv, webm"})
		return
	}

	// Create videos directory if it doesn't exist
	if err := os.MkdirAll("videos", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create videos directory"})
		return
	}

	// Generate unique filename
	filename := filepath.Join("videos", header.Filename)

	// Create the file
	dst, err := os.Create(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Video uploaded successfully",
		"filename": header.Filename,
	})
}

func listVideos(c *gin.Context) {
	files, err := os.ReadDir("videos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read videos directory"})
		return
	}

	var videoFiles []string
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".mkv" || ext == ".webm" {
				videoFiles = append(videoFiles, file.Name())
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"videos": videoFiles})
}
