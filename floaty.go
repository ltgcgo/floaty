// This plugin used luludotdev/caddy-requestid as a starting point.

package floaty

import (
	//"strconv"
	"github.com/caddyserver/caddy/v2"
	//"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	//nanoid "github.com/matoous/go-nanoid/v2"
)

// Floaty implements global placeholders that rolls with a set interval.
type FloatyID struct {
	// Length of instance ID
	Length int `json:"length"`

	// Map of additional instance IDs to be set
	Additional map[string]int `json:"additional,omitempty"`
}

// Initialize the module.
func init() {
	caddy.RegisterModule(FloatyID{})
}

// Register the Caddy plugin.
func (FloatyID) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.system.floaty",
		New: func() caddy.Module {
			return new(FloatyID)
		},
	};
}
