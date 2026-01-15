# goFindMyArm 🕵️‍♂️ 💪

**goFindMyArm** is a tool to easily discover and manage headless miniPCs (Raspberry Pi, Orange Pi, Rock Pi) on your local network.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.23%2B-cyan)
![Wails](https://img.shields.io/badge/wails-v2-red)

## 🚀 Features

-   **Zero-Conf Discovery:** Finds devices via UDP Broadcast (port 32000).
-   **Encrypted Communication:** All payloads are encrypted with NaCl SecretBox.
-   **Remote Commands:** Reboot devices remotely via gRPC (port 32001).
-   **Cross-Platform Agent:** Single binary storage-friendly agent for Linux (ARM/x64).
-   **Modern UI:** Desktop Client built with Vue.js + TailwindCSS.

## 📦 Components

The project is a monorepo containing:

1.  **Agent (`cmd/agent`)**: Runs on the target device (e.g., Raspberry Pi).
2.  **Client (`cmd/client`)**: Desktop application for the user.
3.  **Protocol (`pkg/protocol`)**: Shared Protobuf definitions.

## 🛠️ Build & Install

### Prerequisites
- Go 1.23+
- Node.js 18+
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### 1. Build the Agent
```bash
# For local testing (your current OS)
make agent

# For Rasberry Pi (ARM64)
GOOS=linux GOARCH=arm64 go build -o bin/agent-arm64 ./cmd/agent
```

### 2. Build the Client
```bash
# Run in Dev Mode
make client

# Build for Production
make build-client
```

## 💾 Installation on Device

Run this command on your Raspberry Pi / Linux Device:

```bash
# Opção 1: Link direto (Sempre atualizado com a main)
curl -sL https://raw.githubusercontent.com/apavanello/goFindMyArm/main/scripts/install.sh | sudo bash

# Opção 2: Link Curto (Requer GitHub Pages ativado na pasta /docs)
curl -sL https://apavanello.github.io/goFindMyArm/install | sudo bash
```

**Options:**
-   **Disable Remote Commands:** Append `-- --disable-remote`
-   **Set Password:** Append `-- --password "MY_SECURE_PASS"`

Example:
```bash
curl -sL https://raw.githubusercontent.com/apavanello/goFindMyArm/main/scripts/install.sh | sudo bash -s -- --password "admin123"
```

## 🏗️ Architecture

-   **Discovery:** UDP Broadcast sends a probe. Agents reply with their IP, Hostname, and MAC (Encrypted).
-   **Control:** Client connects via gRPC to Agent IP to execute commands.
-   **Security:** Shared secret (password) used to derive encryption keys and authorize gRPC calls.

## 🤝 Contributing

1.  Fork the repo
2.  Create your feature branch (`git checkout -b feature/amazing-feature`)
3.  Commit your changes (`git commit -m 'Add some amazing feature'`)
4.  Push to the branch (`git push origin feature/amazing-feature`)
5.  Open a Pull Request

## 📄 License

Distributed under the MIT License.
