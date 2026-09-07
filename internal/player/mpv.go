package player

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"path/filepath"
)

func PlayVideo(url string) error {
	ytdlp := exec.Command("yt-dlp", "-o", "-", url, "--js-runtimes", "node", "--remote-components", "ejs:github")
	
	mpvPath := "mpv"
	// Smart fallback for Windows users who use portable MPV
	if runtime.GOOS == "windows" {
		home, _ := os.UserHomeDir()
		portableMpv := filepath.Join(home, "Desktop", "mov", "mpv", "mpv.exe")
		if _, err := os.Stat(portableMpv); err == nil {
			mpvPath = portableMpv
		}
	}

	mpv := exec.Command(mpvPath, "-")

	pipe, err := ytdlp.StdoutPipe()
	if err != nil {
		return err
	}
	mpv.Stdin = pipe

	mpv.Stdout = io.Discard
	mpv.Stderr = io.Discard

	if err := ytdlp.Start(); err != nil {
		return fmt.Errorf("failed to start yt-dlp: %v", err)
	}
	if err := mpv.Start(); err != nil {
		return fmt.Errorf("failed to start mpv: %v", err)
	}

	mpv.Wait()
	ytdlp.Process.Kill()
	return nil
}
