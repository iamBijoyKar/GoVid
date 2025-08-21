# GoVid - Video Streaming Server

A Go-based video streaming server with automatic thumbnail generation and a modern web interface.

## Features

- 🎥 Video streaming with HTTP range requests
- 📤 Video upload support (MP4, AVI, MOV, MKV, WebM)
- 🖼️ Automatic thumbnail generation from video frames
- 🎨 Modern, responsive web interface
- 📱 Grid-based video library with thumbnails
- 🔄 Real-time video list updates

## Prerequisites

- **Go 1.16+** - [Download Go](https://golang.org/dl/)
- **FFmpeg** - Required for thumbnail generation
  - Windows: Download from [FFmpeg.org](https://ffmpeg.org/download.html)
  - macOS: `brew install ffmpeg`
  - Ubuntu/Debian: `sudo apt install ffmpeg`

## Installation

1. Clone or download this repository
2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Start the server
```bash
go run main.go
```

The server will start at `http://localhost:8080`

### Generate thumbnails for existing videos
```bash
go run main.go -generate-thumbnails
```

This will create thumbnails for all videos in the `videos/` directory.

### Upload videos
1. Open `http://localhost:8080` in your browser
2. Use the upload form to add new videos
3. Thumbnails are automatically generated for new uploads

## File Structure

```
GoVid/
├── main.go              # Main server application
├── go.mod               # Go module file
├── go.sum               # Go dependencies checksum
├── static/
│   └── index.html      # Web interface
├── videos/              # Video storage directory
└── thumbnails/         # Generated thumbnails (auto-created)
```

## API Endpoints

- `GET /` - Web interface
- `GET /video?file=<filename>` - Stream video with range support
- `POST /upload` - Upload new video file
- `GET /videos` - List all available videos
- `GET /thumbnail/<filename>` - Get video thumbnail

## Thumbnail Generation

Thumbnails are automatically generated:
- **On upload**: New videos get thumbnails generated in the background
- **On demand**: Thumbnails are created when first requested
- **Batch generation**: Use `-generate-thumbnails` flag for existing videos

Thumbnails are extracted from the 1-second mark of each video for consistency.

## Supported Video Formats

- MP4 (.mp4)
- AVI (.avi)
- MOV (.mov)
- MKV (.mkv)
- WebM (.webm)

## Configuration

- **Port**: Default 8080 (change in `main.go`)
- **File size limit**: 100MB per upload
- **Thumbnail quality**: JPEG with quality 2 (high quality)

## Troubleshooting

### FFmpeg not found
Make sure FFmpeg is installed and available in your system PATH.

### Thumbnails not generating
Check that:
1. FFmpeg is properly installed
2. Video files are valid and readable
3. The `thumbnails/` directory is writable

### Video streaming issues
Ensure your video files are properly encoded and the browser supports the format.

## Development

The application uses:
- **Gin** - HTTP web framework
- **FFmpeg** - Video processing and thumbnail generation
- **Vanilla JavaScript** - Frontend functionality
- **CSS Grid** - Responsive layout

## License

This project is open source and available under the MIT License. 