//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Seed random number generator for varied thumbnail positions
	rand.Seed(time.Now().UnixNano())

	fmt.Println("🎬 GoVid Thumbnail Generator")
	fmt.Println("=============================")
	fmt.Println("🎲 Generating random intro frames (0-10 seconds)")
	fmt.Println()

	// Check if videos directory exists
	if _, err := os.Stat("videos"); os.IsNotExist(err) {
		fmt.Println("❌ Error: 'videos' directory not found!")
		fmt.Println("   Please make sure you're running this script from the GoVid project root.")
		os.Exit(1)
	}

	// Check if ffmpeg is available
	if err := exec.Command("ffmpeg", "-version").Run(); err != nil {
		fmt.Println("❌ Error: FFmpeg not found!")
		fmt.Println("   Please install FFmpeg first:")
		fmt.Println("   - Windows: Download from https://ffmpeg.org/download.html")
		fmt.Println("   - macOS: brew install ffmpeg")
		fmt.Println("   - Ubuntu/Debian: sudo apt install ffmpeg")
		os.Exit(1)
	}

	// Create thumbnails directory if it doesn't exist
	if err := os.MkdirAll("thumbnails", 0755); err != nil {
		fmt.Printf("❌ Failed to create thumbnails directory: %v\n", err)
		os.Exit(1)
	}

	// Read videos directory
	files, err := os.ReadDir("videos")
	if err != nil {
		fmt.Printf("❌ Failed to read videos directory: %v\n", err)
		os.Exit(1)
	}

	// Filter video files
	var videoFiles []string
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			ext := strings.ToLower(filepath.Ext(filename))
			if ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".mkv" || ext == ".webm" {
				videoFiles = append(videoFiles, filename)
			}
		}
	}

	if len(videoFiles) == 0 {
		fmt.Println("ℹ️  No video files found in the 'videos' directory.")
		fmt.Println("   Supported formats: .mp4, .avi, .mov, .mkv, .webm")
		return
	}

	fmt.Printf("📁 Found %d video files\n", len(videoFiles))
	fmt.Println()

	// Process each video file
	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, filename := range videoFiles {
		fmt.Printf("[%d/%d] Processing: %s\n", i+1, len(videoFiles), filename)

		videoPath := filepath.Join("videos", filename)
		thumbnailPath := filepath.Join("thumbnails", strings.TrimSuffix(filename, filepath.Ext(filename))+".jpg")

		// Check if thumbnail already exists
		if _, err := os.Stat(thumbnailPath); err == nil {
			fmt.Printf("   ⏭️  Thumbnail already exists, skipping...\n")
			skipCount++
			continue
		}

		// Generate random timestamp between 0 and 10 seconds for intro section
		randomSeconds := rand.Float64() * 10.0
		timestamp := fmt.Sprintf("%.2f", randomSeconds)

		fmt.Printf("   🎲 Extracting frame at %s seconds...\n", timestamp)

		// Generate thumbnail using ffmpeg with random timestamp
		startTime := time.Now()
		cmd := exec.Command("ffmpeg",
			"-i", videoPath,
			"-ss", timestamp,
			"-vframes", "1",
			"-q:v", "2",
			"-y", // Overwrite output files
			thumbnailPath)

		// Capture ffmpeg output for better error reporting
		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("   ❌ Failed to generate thumbnail: %v\n", err)
			if len(output) > 0 {
				fmt.Printf("   FFmpeg output: %s\n", string(output))
			}
			errorCount++
		} else {
			duration := time.Since(startTime)
			fmt.Printf("   ✅ Generated thumbnail in %v\n", duration.Round(time.Millisecond))
			successCount++
		}

		fmt.Println()
	}

	// Summary
	fmt.Println("🎯 Thumbnail Generation Complete!")
	fmt.Println("================================")
	fmt.Printf("✅ Successfully generated: %d\n", successCount)
	fmt.Printf("⏭️  Skipped (already exist): %d\n", skipCount)
	fmt.Printf("❌ Failed: %d\n", errorCount)
	fmt.Printf("📊 Total processed: %d\n", len(videoFiles))

	if errorCount > 0 {
		fmt.Println("\n⚠️  Some thumbnails failed to generate.")
		fmt.Println("   Check the error messages above for details.")
		os.Exit(1)
	} else {
		fmt.Println("\n🎉 All thumbnails generated successfully!")
		fmt.Println("🎲 Each thumbnail shows a random frame from the intro section!")
	}
}
