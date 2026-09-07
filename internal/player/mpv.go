package player

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func PlayVideo(url string) error {
	ytdlp := exec.Command("yt-dlp", "-o", "-", url, "--js-runtimes", "node", "--remote-components", "ejs:github")
	
	mpvPath := os.ExpandEnv("C:\\Users\\lross\\Desktop\\mov\\mpv\\mpv.exe")
	if _, err := os.Stat(mpvPath); os.IsNotExist(err) {
		mpvPath = "mpv"
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
