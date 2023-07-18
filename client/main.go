package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/gen2brain/beeep"
)

func captureScreenshotCmd(outputFileName string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("gnome-screenshot", "--file", outputFileName, "--area")
	case "darwin":
		cmd = exec.Command("screencapture", "-i", outputFileName)
	case "windows":
		cmd = exec.Command("snippingtool", "/clip")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error capturing screenshot: %s", err)
	}

	return nil
}

func uploadImage(filename string, uploadURL string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	apiKeyPart, err := writer.CreateFormField("api_key")
	if err != nil {
		return err
	}

	_, err = apiKeyPart.Write([]byte("deez_nuts")) //replace your key here bruh

	if err != nil {
		return err
	}

	filePart, err := writer.CreateFormFile("screenshot", filepath.Base(filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(filePart, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload image, status code: %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if buf.String() == "Invalid API key" {
		err := fmt.Errorf("invalid API key")
		return err
	}

	clipboard.WriteAll(buf.String())

	return nil
}

func generateRandomHash() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Error generating random hash: ", err)
	}
	return hex.EncodeToString(bytes)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: deezshots-cleint <Upload URL>")
		return
	}

	tmpFileName := generateRandomHash() + ".png"

	err := captureScreenshotCmd(tmpFileName)
	if err != nil {
		beeep.Alert("Error", "Error capturing screenshot", "")
		return
	}

	uploadURL := os.Args[1]
	if err := uploadImage(tmpFileName, uploadURL); err != nil {
		if err.Error() == "invalid API key" {
			beeep.Alert("Error", "Invalid API key", "")
			return
		}

		beeep.Alert("Error", "Error uploading image", "")
		os.Remove(tmpFileName)
		return
	}

	err = os.Remove(tmpFileName)

	if err != nil {
		beeep.Alert("Error", "Error removing tmp file", "")
		return
	}

	beeep.Notify("Success", "Screenshot uploaded successfully", "")
}
