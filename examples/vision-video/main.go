package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jonathanhecl/gollama"
)

const (
	pathScreenshots = "./screenshots/"
)

var (
	videoExt = []string{".mp4", ".mkv", ".webm"}
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	m := gollama.New("llama3.2-vision")
	m.PullIfMissing(ctx)

	// Path with videos
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Insert path with videos:")
	if !scanner.Scan() {
		fmt.Println("Error reading path")
		return
	}
	path := scanner.Text()

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Path does not exist")
		return
	}

	// List videos on the path
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading path")
		return
	}

	type sVideo struct {
		path     string
		filename string
	}

	videos := []sVideo{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())

		for _, v := range videoExt {
			if ext == v {
				videos = append(videos, sVideo{path: path, filename: file.Name()})
				break
			}
		}
	}

	// Print videos
	fmt.Println("Videos:")
	for _, video := range videos {
		fmt.Println(video.filename)

		// Extract first frame
		extractFrame(video.path, video.filename)
	}

	fmt.Println("Total videos:", len(videos))

	// Request a prompt
	fmt.Println("Insert request:")
	if !scanner.Scan() {
		fmt.Println("Error reading prompt")
		return
	}
	prompt := scanner.Text()

	fmt.Println("Prompt:", prompt)

	for _, video := range videos {
		fmt.Println("Processing video:", video.filename)

		image := gollama.PromptImage{Filename: pathScreenshots + video.filename + ".jpg"}

		res, err := m.Chat(ctx, prompt, image)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Response:", res)
	}
}

func extractFrame(path string, filename string) {
	// Check if pathScreenshots exists
	if _, err := os.Stat(pathScreenshots); os.IsNotExist(err) {
		os.Mkdir(pathScreenshots, 0755)
	}

	command := "ffmpeg"
	frameExtractionTime := "0:00:05.000"
	output := pathScreenshots + filename + ".jpg"

	if _, err := os.Stat(output); !os.IsNotExist(err) {
		return // Skip if output file already exists
	}

	cmd := exec.Command(command,
		"-ss", frameExtractionTime,
		"-i", filepath.Join(path, filename),
		output)

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()
}
