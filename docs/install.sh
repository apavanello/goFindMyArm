#!/bin/bash
set -e

# Defaults
REPO="apavanello/goFindMyArm"
INSTALL_DIR="/usr/local/bin"
SERVICE_NAME="gofindmyarm"
BINARY_NAME="gofindmyarm-agent"
QUIET=false
DISABLE_REMOTE=false
PASSWORD=""

# Print usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  --quiet               Non-interactive mode (auto-generates password if missing)"
    echo "  --disable-remote      Disable remote commands (Reboot/Update)"
    echo "  --password STR        Set password for remote commands"
    echo "  --help                Show this help"
    exit 1
}

# Parse Args
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --quiet) QUIET=true ;;
        --disable-remote) DISABLE_REMOTE=true ;;
        --password) PASSWORD="$2"; shift ;;
        --help) usage ;;
        *) echo "Unknown parameter: $1"; usage ;;
    esac
    shift
done

# Root check
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

echo ">>> goFindMyArm Installer"

# 1. Architecture Detection
ARCH=$(uname -m)
case $ARCH in
    x86_64) GOARCH="amd64" ;;
    aarch64) GOARCH="arm64" ;;
    armv7l) GOARCH="arm" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac
echo "Detected Architecture: $GOARCH"

# 2. Password Logic
if [ "$DISABLE_REMOTE" = true ]; then
    echo "Remote commands DISABLED by user request."
elif [ -z "$PASSWORD" ]; then
    if [ "$QUIET" = true ]; then
        # Generate random password
        PASSWORD=$(openssl rand -hex 12)
        echo "Generated Password: $PASSWORD"
        echo "SAVE THIS PASSWORD!"
    else
        # Interactive prompt
        echo -n "Enter password for remote commands: "
        # Fix for curl | bash: read from /dev/tty
        if [ -t 0 ]; then
            read -s PASSWORD
        else
            if [ -c /dev/tty ]; then
                read -s PASSWORD < /dev/tty
            else
                echo "Error: Cannot read password (no tty). Use --password or --quiet."
                exit 1
            fi
        fi
        echo
        if [ -z "$PASSWORD" ]; then
             echo "Password cannot be empty!"
             exit 1
        fi
    fi
fi

# 3. Download Binary
# Fetch latest tag from GitHub API
echo "Fetching latest version info..."
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")' || true)

if [ -z "$LATEST_TAG" ]; then
    echo "Error: Could not find latest release on GitHub."
    echo "Check if the repository has any releases at https://github.com/$REPO/releases"
    exit 1
fi

echo "Latest Version: $LATEST_TAG"

echo "--------------------------------------------------"
echo "Installing goFindMyArm Agent ($LATEST_TAG)"
echo "--------------------------------------------------"

# Construct URL (e.g. agent-linux-arm64)
BINARY_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/agent-linux-$GOARCH"

echo "Downloading from $BINARY_URL..."
curl -L -o "$INSTALL_DIR/$BINARY_NAME" "$BINARY_URL" 

if [ $? -ne 0 ]; then
    echo "Error: Failed to download binary. Check architecture ($GOARCH) or network."
    exit 1
fi

chmod +x "$INSTALL_DIR/$BINARY_NAME" || true

# 4. Systemd Service
echo "Creating Systemd Service..."
cat <<EOF > /etc/systemd/system/$SERVICE_NAME.service
[Unit]
Description=goFindMyArm Agent Service
After=network.target

[Service]
Type=simple
ExecStart=$INSTALL_DIR/$BINARY_NAME --udp 32000 --tcp 32001 --disable-remote=$DISABLE_REMOTE --password="$PASSWORD"
Restart=always
RestartSec=5
User=root

[Install]
WantedBy=multi-user.target
EOF

# 5. Enable & Start
systemctl daemon-reload
# systemctl enable --now $SERVICE_NAME (Uncomment to actually start)

echo ">>> Installation Complete!"
if [ "$DISABLE_REMOTE" = false ]; then
    echo "Remote commands enabled."
fi
echo "Run 'systemctl start $SERVICE_NAME' to start."
