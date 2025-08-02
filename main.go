package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.StaticFile("/", "./static")
	r.GET("/video", streamVideo)
	fmt.Println("Server started at http://localhost:8080/video")
	r.Run(":8080")
}

func streamVideo(c *gin.Context) {
	filePath := "videos/sample.mp4"
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
