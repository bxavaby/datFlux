<img src="assets/df-tn.png" alt="datFlux Logo" width="200"/>

> An entropy-borne password generator with a Tokyo Night TUI dashboard

datFlux is a terminal-based password generator that uses system noise as entropy sources to create truly random passwords. It features a beautiful, responsive terminal interface and creates cryptographically secure passwords while showing real-time system metrics.

<br>

## ☍ Screenshots

<p align="left">
  <img src="assets/df-animate.png" alt="datFlux anim" width="500"><br>
  <em>⊹ mid-animation ⊹</em>
</p>

<p align="left">
  <img src="assets/df-pwd.png" alt="datFlux gen" width="500"><br>
  <em>⊹ password generated ⊹</em>
</p>

<p align="left">
  <img src="assets/df-copy.png" alt="datFlux cpy" width="500"><br>
  <em>⊹ copy password ⊹</em>
</p>

<br>

## ☍ System Requirements

- **Linux** (primary support): full functionality available
- **macOS** (partial support): basic features work, but system monitoring may be limited
- **Windows**: not currently supported

For best results, try datFlux on a Linux system.

<br>

## ☍ Features

- **High-Entropy Password Generation**: creates strong passwords using real system noise
- **Cinematic Password Reveal Animation**: visually decrypts passwords character by character
- **System Metrics Dashboard**: live CPU, RAM, and Network usage monitoring
- **Entropy Quality Indicator**: monitors randomness quality in real-time
- **Clipboard Integration**: copy passwords to your clipboard with a single keystroke

<br>

## ☍ Installation

### Quick Install Script (Linux/macOS)

The easiest way to install datFlux is using the installation script:

```bash
# Clone the repository
git clone https://github.com/bxavaby/datFlux.git
cd datFlux

# Run the installer
./install.sh
```

The installer will build datFlux and add it to your system path.

### From Source (Manual)

```bash
# Clone the repository
git clone https://github.com/bxavaby/datFlux.git
cd datFlux

# Build the binary
go build -o datflux ./cmd/datflux

# Run it
./datflux
```

### Go Install

```bash
go install github.com/bxavaby/datFlux/cmd/datflux@latest
```

<br>

## ☍ Usage

Launch it in your terminal:

```bash
datflux
```

### Key Commands

- <kbd>r</kbd> - (re)generate password
- <kbd>c</kbd> - copy the generated password to clipboard
- <kbd>q</kbd> / <kbd>Ctrl+C</kbd> / <kbd>Esc</kbd> - quit datFlux

<br>

## ☍ How It Works

datFlux creates background system load through various noise generation methods:

1. **CPU Noise**: performs complex mathematical operations
2. **RAM Noise**: allocates and manipulates memory blocks
3. **Network Noise**: creates local network connections and data transfer

These operations generate entropy that is collected, hashed, and used to create unpredictable, secure passwords that are more resistant to brute force and dictionary attacks than traditional password generators.

<br>

## ☍ Security Considerations

- datFlux generates passwords locally; no data is sent over the network
- The program creates system load to generate entropy but has safeguards to prevent excessive resource usage
- Generated passwords never leave your computer unless you copy them

<br>

## ☍ License

MIT License - see [LICENSE](LICENSE) for details

---

<p align="center">By <a href="https://github.com/bxavaby">bxavaby</a></p>
