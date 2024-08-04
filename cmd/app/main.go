package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request: %v", err)
			return
		}

		// Save the uploaded file to disk
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "./upload/"+filename); err != nil {
			c.String(http.StatusInternalServerError, "Could not save file: %v", err)
			return
		}

		// Create output directory
		hlsDir := "./output"
		os.MkdirAll(hlsDir, os.ModePerm)

		// Start video conversion in a goroutine
		go convertVideo("./upload/"+filename, hlsDir)

		c.String(http.StatusOK, "Video conversion started")
	})

	r.Run(":8889")
}

func convertVideo(filename, hlsDir string) {
	// Resolutions to be used for HLS
	resolutions := []struct {
		Width   int
		Height  int
		Bitrate string
	}{
		{640, 360, "800k"},
		{1280, 720, "2800k"},
		{1920, 1080, "5000k"},
	}

	// Generate HLS for each resolution
	for _, res := range resolutions {
		resDir := filepath.Join(hlsDir, fmt.Sprintf("%dx%d", res.Width, res.Height))
		os.MkdirAll(resDir, os.ModePerm)
		hlsPath := filepath.Join(resDir, "index.m3u8")

		cmd := exec.Command("ffmpeg", "-i", filename, "-vf", fmt.Sprintf("scale=%d:%d", res.Width, res.Height),
			"-c:a", "aac", "-ar", "48000", "-b:a", "128k", "-c:v", "h264", "-profile:v", "main", "-crf", "20",
			"-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-hls_time", "4", "-hls_playlist_type", "vod",
			"-b:v", res.Bitrate, "-maxrate", res.Bitrate, "-bufsize", "1000k", "-hls_segment_filename",
			filepath.Join(resDir, "segment_%03d.ts"), hlsPath)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to convert video to resolution %dx%d: %v\n", res.Width, res.Height, err)
			return
		}
	}

	// Create master playlist
	masterPlaylist := filepath.Join(hlsDir, "master.m3u8")
	masterFile, err := os.Create(masterPlaylist)
	if err != nil {
		fmt.Printf("Failed to create master playlist: %v\n", err)
		return
	}
	defer masterFile.Close()

	masterFile.WriteString("#EXTM3U\n")
	for _, res := range resolutions {
		resDir := filepath.Join(fmt.Sprintf("%dx%d", res.Width, res.Height))
		hlsPath := filepath.Join(resDir, "index.m3u8")
		masterFile.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%s,RESOLUTION=%dx%d\n%s\n",
			res.Bitrate, res.Width, res.Height, hlsPath))
	}

	fmt.Printf("Video converted successfully: %s\n", masterPlaylist)
}
