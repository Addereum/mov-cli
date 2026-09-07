# mov-cli (V5 Alpha)

The next generation of `mov-cli`, written entirely in Go.

## Why Go?
- **Single Binary**: No Python, no dependencies, no hassle.
- **Blazing Fast**: Native execution with Goroutine support.
- **JavaScript Plugin Engine**: Write lightweight scrapers in pure JS that run on our embedded Go VM!

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
*(Note: During this Alpha phase, ensure the `plugins/` folder remains in your working directory so the Go engine can load the JavaScript scrapers).*
