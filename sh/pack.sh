#!/bin/bash
# shx Standard Utility
#
# Customizable bundle emitter (upcoming)
#
# Clean up the previous builds
mkdir -p build
rm build/release.* 2>/dev/null
echob "Creating TAR archive..."
tar --verbose -cf build/release.tar \
	shx .gitignore \
	nix/nano nix/env.nix nix/zsh/.zshName nix/zsh/.zshrc \
	sh/sh.sh sh/commit.sh sh/amend.sh sh/tag.sh sh/push.sh \
	sh/which sh/echob\
	LICENSE README.md
echob "Emitting GZIP archive..."
if [ -f "$(which zopfli)" ]; then
	zopfli -v --i25 build/release.tar 2>&1 | while IFS= read -r line; do
		if [[ "$line" == "block split"* ]]; then
			echo $line
		elif [[ "$line" == "Original "* ]]; then
			echo $line
		fi
	done
	mv build/release.tar.gz build/release.tgz
elif [ -f $(which gzip) ]; then
	gzip -kvc9 build/release.tar > build/release.tgz
else
	echob "No GZIP-compatible compression utility found. Skipped."
fi
echob "Emitting Brotli archive..."
if [ -f "$(which brotli)" ]; then
	brotli -kvq 11 build/release.tar
	mv build/release.tar.br build/release.tbr
else
	echob "No Brotli compression utility found. Skipped."
fi