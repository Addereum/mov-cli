# mov-cli (V5 Alpha)

The next generation of `mov-cli`, written entirely in Go.

> ⚠️ **Origins & Credits**: This project is a modern GoLang rewrite inspired by the original Python-based [mov-cli](https://github.com/mov-cli/mov-cli). Because the original developers have abandoned the project and it is no longer maintained, this repository was created to keep the project alive. We rebuilt it from the ground up in Go to solve previous dependency issues and make it faster!

## Why Go?
- **Single Binary**: No Python, no dependencies, no hassle.
- **Blazing Fast**: Native execution with Goroutine support.
- **JavaScript Plugin Engine**: Write lightweight scrapers in pure JS that run on our embedded Go VM!

## 🧩 Plugin Ecosystem (New in V5)
V5 comes with a built-in Plugin Manager. Plugins are simply lightweight `.js` files that run securely in our embedded JavaScript engine.

**Installing a plugin via URL (1-Click):**
```bash
mov-cli plugin add https://raw.githubusercontent.com/Addereum/mov-cli/v5/plugins/youtube.js
```
The CLI automatically downloads the script to `~/.config/mov-cli/plugins/` and makes it immediately available.

### 📚 Plugin Developer Documentation
Want to write your own scraper? It's incredibly easy. We provide native HTTP fetching and HTML parsing (like jQuery) directly inside the JavaScript environment, meaning you don't even need `yt-dlp` or external libraries for most sites!

* **Read the [Plugin API Documentation](docs/plugin_api.md)**
* **Check out the [Example Scraper](examples/dummy_scraper.js)**

## 🐧 Linux / Arch Installation

### Arch Linux (AUR)
We provide a native `PKGBUILD` for Arch Linux users that compiles the Go binary and sets everything up.
```bash
git clone -b v5 https://github.com/Addereum/mov-cli.git
cd mov-cli/aur
makepkg -si
```

### Universal Linux / macOS (Go Install)
If you have Go installed on your system, you can pull and install the V5 branch directly:
```bash
go install github.com/Addereum/mov-cli@v5
```

### Manual Build
```bash
git clone -b v5 https://github.com/Addereum/mov-cli.git
cd mov-cli
go build -o mov-cli
./mov-cli "Spongebob"
```
