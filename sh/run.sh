#!/bin/bash
xcaddy run 2>&1 | while IFS= read -r line; do
	echo "${line//${SOURCE_DIR}/\~nix}"
done
exit