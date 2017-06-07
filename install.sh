#!/bin/sh

TID_ARCH=""
TID_OS=""

TID_DIR="$HOME/.tid"
TID_BIN_DIR="$TID_DIR/bin"
TID_VERSION=""

prepareArch() {
	TID_ARCH=$(uname -m)

	case "$TID_ARCH" in
		x86) TID_ARCH="386";;
		x86_64) TID_ARCH="amd64";;
		i686) TID_ARCH="386";;
		i386) TID_ARCH="386";;
	esac
}

prepareOS() {
	TID_OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

	case "$TID_OS" in
		# Minimalist GNU for Windows
		mingw*) TID_OS='windows';;
	esac
}

prepareTidDir() {
    if test ! -d "$TID_DIR"; then
        mkdir -p "$TID_DIR"
    fi

	if test ! -d "$TID_BIN_DIR"; then
		mkdir -p "$TID_BIN_DIR"
	fi
}

prepareVersion() {
    if type "curl" > /dev/null; then
		TID_VERSION=$(curl -s https://raw.githubusercontent.com/SeerUK/tid/master/VERSION)
	elif type "wget" > /dev/null; then
		TID_VERSION=$(wget -q -O - https://raw.githubusercontent.com/SeerUK/tid/master/VERSION)
	fi
}

installRelease() {
	if test -z "$TID_ARCH"; then
		echo "Sorry, we it looks like your system architecture is unsupported."
		echo "You can always compile the Go source from https://github.com/SeerUK/tid though!"
		exit 1
	fi

	if test -z "$TID_OS"; then
		echo "Sorry, we it looks like your operating system is unsupported."
		echo "You can always compile the Go source from https://github.com/SeerUK/tid though!"
		exit 1
	fi

	BINARIES="tid tiddate"

	for BINARY in ${BINARIES}; do
		BINARY_RELEASE_URL="https://github.com/SeerUK/tid/releases/download/${TID_VERSION}/${BINARY}_${TID_OS}_${TID_ARCH}"

		echo "Downloading $BINARY $TID_VERSION..."

		BINARY_VERSION_DIR="${TID_BIN_DIR}-${TID_VERSION}"
		BINARY_VERSION_FILE="$BINARY_VERSION_DIR/$BINARY"

		if test "$TID_OS" = "windows"; then
			BINARY_VERSION_FILE="$BINARY_VERSION_FILE.exe"
		fi

		mkdir -p "$BINARY_VERSION_DIR"

		if type "curl" > /dev/null; then
			curl -L "$BINARY_RELEASE_URL" -o "$BINARY_VERSION_FILE"
		elif type "wget" > /dev/null; then
			wget -q -O "$BINARY_VERSION_FILE" "$BINARY_RELEASE_URL"
		fi

		if test "$TID_OS" != "windows"; then
		    chmod +x "$BINARY_VERSION_FILE"
		fi

		echo "Linking $BINARY $TID_VERSION..."

		ln -fs "$BINARY_VERSION_FILE" "$TID_BIN_DIR/${BINARY}"
	done

	echo ""
	echo "Installed to '$TID_BIN_DIR'."
	echo "This should be have automatically been added to your path."
	echo "If it hasn't been, make sure you add it manually!"
}

installPath() {
    BASH_LINE="export PATH=\"\$PATH:\$HOME/.tid/bin"
    FISH_LINE="set -gx PATH \$HOME/.tid/bin \$PATH"
    ZSH_LINE="$BASH_LINE"

    grep -qF "$BASH_LINE" "$HOME/.bashrc" ||
        echo "$BASH_LINE" >> "$HOME/.bashrc"

    grep -qF "$BASH_LINE" "$HOME/.profile" ||
        echo "$BASH_LINE" >> "$HOME/.profile"

    grep -qF "$FISH_LINE" "$HOME/.config/fish/config.fish" ||
        echo "$FISH_LINE" >> "$HOME/.config/fish/config.fish"

    grep -qF "$ZSH_LINE" "$HOME/.zshrc" ||
        echo "$ZSH_LINE" >> "$HOME/.zshrc"
}

prepareArch
prepareOS
prepareTidDir
prepareVersion

installRelease
installPath

# @todo:
# Backup existing data (how do we know where the DB is? parse config for path?)
