# Video to HLS Converter Service
![Video to HLS Converter Service](./assets/anime.gif)

This Go-based service allows you to upload video files and converts them to HLS (HTTP Live Streaming) format with multiple resolutions. The service uses `ffmpeg` to handle video conversion and runs the conversion process in a goroutine for asynchronous processing.

## Features

- Upload video files via HTTP POST request.
- Convert video files to HLS format with multiple resolutions (360p, 720p, 1080p).
- Generate a master playlist (`master.m3u8`) that includes all resolutions.
- Asynchronous video conversion using goroutines.

## Requirements

- Go 1.16 or higher
- `ffmpeg` installed on your system

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/digkill/HLS-video-converter video-to-hls
   cd video-to-hls
