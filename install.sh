#!/bin/bash


install() {
	set -eu
	UNAME=$(uname)

	if [ "$UNAME" != "Linux" -a "$UNAME" != "Darwin" ] ; then
		echo "Sorry, OS not supported: ${UNAME}."
		exit 1
	fi

	if [ "$UNAME" = "Darwin" ] ; then
		OSX_ARCH=$(uname -m)
		if [ "${OSX_ARCH}" = "x86_64" ] ; then
			PLATFORM="darwin-amd64"
		else
			echo "Sorry, architecture not supported: ${OSX_ARCH}."
			exit 1
		fi
	elif [ "$UNAME" = "Linux" ] ; then
		LINUX_ARCH=$(uname -m)
		if [ "${LINUX_ARCH}" = "x86_64" ] ; then
			PLATFORM="linux-amd64"
		else
			echo "Sorry, architecture not supported: ${LINUX_ARCH}."
			exit 1
		fi
	fi
	LATEST=$(curl -s https://api.github.com/repos/jobtalk/pnzr/tags | grep -Eo '"name":.*[^\\]",'  | head -n 1 | sed 's/[," ]//g' | cut -d ':' -f 2)
	URL="https://github.com/jobtalk/pnzr/releases/download/$LATEST/pnzr-$PLATFORM"
	DEST=${DEST:-/usr/local/bin/pnzr}

	if [ -z $LATEST ] ; then
		echo "Error requesting. Download binary from https://github.com/jobtalk/pnzr/releases"
		exit 1
	else
		echo "Downloading pnzr binary from https://github.com/jobtalk/pnzr/releases/download/$LATEST/pnzr-$PLATFORM to $DEST"
		if curl -sL https://github.com/jobtalk/pnzr/releases/download/$LATEST/pnzr-$PLATFORM -o $DEST; then
			chmod +x $DEST
			echo "pnzr installation was successful"
		else
			echo "Installation failed. You may need elevated permissions."
			exit 1
		fi
	fi
}

install
