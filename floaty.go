// This plugin used luludotdev/caddy-requestid as a starting point.

package floaty

import (
	"net/http"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	nanoid "github.com/matoous/go-nanoid/v2"
)

// Floaty implements global placeholders that rolls with a set interval.
type FloatyID struct {
	// Length of instance ID
	Length int `json:"length"`;

	// Map of additional instance IDs to be set
	Additional map[string]int `json:"additional,omitempty"`;
}

// Initialize the module.
