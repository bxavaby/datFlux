<div align="center">
  <img src="assets/df-lg.png" alt="datFlux Logo" width="300"/>

  <h1>datFlux</h1>
  <p><em>An entropy-borne password generator with a Tokyo Night TUI dashboard</em></p>

  <a href="#-features"><img src="https://img.shields.io/badge/Features-Overview-1a1b26?color=7aa2f7" alt="Features"></a>
  <a href="#-installation"><img src="https://img.shields.io/badge/Installation-Guide-1a1b26?color=9ece6a" alt="Installation"></a>
  <a href="#-usage"><img src="https://img.shields.io/badge/Usage-Instructions-1a1b26?color=ff9e64" alt="Usage"></a>
  <a href="#-security-considerations"><img src="https://img.shields.io/badge/Security-Considerations-1a1b26?color=bb9af7" alt="Security"></a>

  <br>
</div>

<div align="center">
  <a href="https://github.com/yourusername/datflux">
    <img src="https://readme-typing-svg.herokuapp.com?font=JetBrains+Mono&size=18&duration=3500&pause=100&color=bb9af7&center=true&vCenter=true&multiline=true&width=800&height=120&lines=datFlux+is+a+terminal-based+password+generator+that+uses+system+noise;as+entropy+sources+to+create+truly+random+passwords+-+it+features;a+beautiful%2C+responsive+TUI+and+creates+cryptographically;secure+passwords+while+showing+real-time+system+metrics" alt="datFlux Overview"/>
  </a>
  <br><br>
</div>

<div align="center">

## § Walkthrough

<br>

</div>

<div align="center">
  <h3>⊹ Password Generation ⊹</h3>
  <p><em>Mid-animation password reveal</em></p>
  <img src="assets/df-animate.png" alt="datFlux animation" width="500">
  <br><br>

  <h3>⊹ Copy to Clipboard ⊹</h3>
  <p><em>One-key operation to copy the password</em></p>
  <img src="assets/df-copy.png" alt="Copy password" width="500">
  <br><br>

  <h3>⊹ Strength Analysis ⊹</h3>
  <p><em>Breakdown of the password's security</em></p>
  <img src="assets/df-pwdstg.png" alt="Password strength" width="500">
</div>

<br>

#

<div align="center">
  <p>datFlux comes with multiple colorschemes to match your aesthetic:</p>
  <br>

  <h3>Tokyo Night (Default)</h3>
  <p><em>Inspired by the beautiful city of Tokyo at night</em></p>
  <a href="https://github.com/enkia/tokyo-night-vscode-theme" target="_blank" rel="noopener noreferrer">
    <img src="assets/df-tn.png" alt="Tokyo Night theme" width="500">
  </a>
  <p><small>Palette by <a href="https://github.com/enkia" target="_blank" rel="noopener noreferrer">enkia</a></small></p>
  <br><br>

  <h3>Ozone-10</h3>
  <p><em>Inspired by polluted city sky colors</em></p>
  <a href="https://lospec.com/palette-list/ozone-10" target="_blank" rel="noopener noreferrer">
    <img src="assets/df-ozone10.png" alt="Ozone-10 theme" width="500">
  </a>
  <p><small>Palette by <a href="https://lospec.com/tinapxl" target="_blank" rel="noopener noreferrer">@tinapxl</a></small></p>
  <br><br>

  <h3>Hydrangea 11</h3>
  <p><em>Inspired by hydrangea flowers</em></p>
  <a href="https://lospec.com/palette-list/hydrangea-11" target="_blank" rel="noopener noreferrer">
    <img src="assets/df-hydrangea11.png" alt="Hydrangea 11 theme" width="500">
  </a>
  <p><small>Palette by <a href="https://lospec.com/dinchenix" target="_blank" rel="noopener noreferrer">@dinchenix</a></small></p>
  <br><br>

  <h3>Leopold's Dreams</h3>
  <p><em>Inspired by blueish melancholic sci-fi scenarios</em></p>
  <a href="https://lospec.com/palette-list/leopolds-dreams" target="_blank" rel="noopener noreferrer">
    <img src="assets/df-leopoldsdreams.png" alt="Leopold's Dreams theme" width="500">
  </a>
  <p><small>Palette by <a href="https://lospec.com/sukinapan" target="_blank" rel="noopener noreferrer">@sukinapan</a></small></p>
</div>

<br><br>

> ※ Press <kbd>t</kbd> anytime to cycle between available themes. Your preference will be applied immediately without interrupting the workflow.

<br>

<div align="center">

## § System Requirements

![Linux: Full Support](https://img.shields.io/badge/Linux-Full%20Support-success?logo=linux&logoColor=white)
![macOS: Partial Support](https://img.shields.io/badge/macOS-Partial%20Support-yellow?logo=apple&logoColor=white)
![Windows: Not Supported](https://img.shields.io/badge/Windows-Not%20Supported-critical?logo=windows&logoColor=white)

<br>

<p>
  <strong>Linux</strong>: full functionality available<br>
  <strong>macOS</strong>: system monitoring may be limited<br>
  <strong>Windows</strong>: not currently supported
</p>

</div>

<br><br>

> ※ For the best experience, try datFlux on a Linux system.

<br>

<div align="center">

## § Features

  <p><strong>High-Entropy Password Generation</strong><br>creates strong passwords using real system noise</p>
  <br>

  <p><strong>Cinematic Password Reveal Animation</strong><br>visually decrypts passwords character by character</p>
  <br>

  <p><strong>Multiple Attack Models</strong><br>test against different threat scenarios from online attacks to QC</p>
  <br>

  <p><strong>System Metrics Dashboard</strong><br>live CPU, RAM, and Network usage monitoring</p>
  <br>

  <p><strong>Entropy Quality Indicator</strong><br>monitors randomness quality in real-time</p>
  <br>

  <p><strong>Clipboard Integration</strong><br>copy passwords to your clipboard with a single keystroke</p>
</div>

<br><br>

> ※ The quantum computing attack model provides time estimates with approximately 5-8% margin of error when compared to theoretical calculations. These estimates are rounded to the nearest time unit for readability.

<br><br>

<div align="center">

## § Installation

  <h3>Quick Install Script (Linux/macOS)</h3>
  <p>The easiest way to install datFlux is using the installation script:</p>

```bash
# Clone the repository
git clone https://github.com/bxavaby/datFlux.git
cd datFlux

# Run the installer
./install.sh
```

  <p>The installer will build datFlux and add it to your system path.</p>
  <br>

  <h3>From Source (Manual)</h3>

```bash
# Clone the repository
git clone https://github.com/bxavaby/datFlux.git
cd datFlux

# Build the binary
go build -o datflux ./cmd/datflux

# Run it
./datflux
```

<br>

  <h3>Go Install</h3>

```bash
go install github.com/bxavaby/datFlux/cmd/datflux@latest
```
</div>

<br><br>

<div align="center">

## § Usage

  <p>Launch it in your terminal:</p>

```bash
datflux
```

  <h3>Key Commands</h3>

  <p>
    <kbd>r</kbd> - generate password<br>
    <kbd>c</kbd> - copy the password<br>
    <kbd>o</kbd> - cycle attack models<br>
    <kbd>t</kbd> - cycle through themes<br>
    <kbd>q</kbd> / <kbd>Ctrl+C</kbd> / <kbd>Esc</kbd> - quit datFlux
  </p>
</div>

<br><br>

<div align="center">

## § How It Works

  <p>datFlux creates background system load through various noise generation methods:</p>
  <br>

  <p><strong>1. CPU Noise</strong><br>performs complex mathematical operations</p>
  <br>

  <p><strong>2. RAM Noise</strong><br>allocates and manipulates memory blocks</p>
  <br>

  <p><strong>3. Network Noise</strong><br>creates local network connections and data transfer</p>
  <br>

  <p>These operations generate entropy that is collected, hashed, and used to create unpredictable, secure passwords that are more resistant to brute force and dictionary attacks than traditional password generators.</p>
</div>

<br><br>

<div align="center">

## § Security Considerations

  <p>• datFlux password generation happens locally with zero network transmission</p>
  <p>• System load is optimized with safeguards to collect entropy efficiently</p>
  <p>• Your passwords remain on-device until you explicitly copy them elsewhere</p>
</div>

<br><br>

#

<div align="center">
  <p><em>v1.0.0 — Entropy-driven security</em></p>
</div>

<div align="center">
  <p>
    <a href="https://github.com/bxavaby" target="_blank" rel="noopener noreferrer">
      <img src="https://img.shields.io/badge/crafted_with_%E2%9D%A4%EF%B8%8F_by-bxavaby-7aa2f7?style=for-the-badge" alt="Created by bxavaby">
    </a>
  </p>
  <p><a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-bb9af7?style=for-the-badge" alt="MIT License"></a></p>
</div>
