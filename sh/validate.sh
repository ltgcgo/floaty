#!/bin/bash
xcaddy build --with github.com/lolPants/caddy-requestid --with github.com/porech/caddy-maxmind-geolocation 2>&1 | while IFS= read -r line; do
	echo "${line//${SOURCE_DIR}/\~nix}"
done
./caddy run -c ./Caddyfile.validate
exit