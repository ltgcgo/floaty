// This plugin used luludotdev/caddy-requestid as a starting point

package floaty

import (
	//"strconv"
	"github.com/caddyserver/caddy/v2"
	//"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	nanoid "github.com/matoous/go-nanoid/v2"
)

// Initialization phase

// Floaty implements global placeholders that rolls with a set interval
type FloatyID struct {
	// Length of instance ID
	Length int `json:"length"`

	// Map of additional instance IDs to be set
	Additional map[string]int `json:"additional,omitempty"`

	// Initialized master ID
	InstanceId string

	// Map of additional instance IDs initialized
	MappedIds map[string]string
}

// Initialize the module
func init() {
	caddy.RegisterModule(FloatyID{});
}

// Register the Caddy plugin
func (FloatyID) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.floaty",
		New: func() caddy.Module {
			return new(FloatyID)
		},
	};
}

// Provisioning phase

// Set up the IDs
var floatyIdGlobal string;
var floatyIdMapped map[string]string;
func (m *FloatyID) Provision(ctx caddy.Context) error {
	// Normalize the parameters
	if m.Length < 1 {
		m.Length = 8;
	};
	if m.Additional == nil {
		m.Additional = make(map[string]int);
	};
	// Generate the IDs
	floatyIdGlobal = nanoid.Must(m.Length);
	for i0, e0 := range m.Additional {
		if e0 < 1 {
			e0 = 8;
		};
		floatyIdMapped[i0] = nanoid.Must(e0);
	};
	// Bind the IDs to global variables
	m.InstanceId = floatyIdGlobal;
	m.MappedIds = floatyIdMapped;
	return nil;
}

// Handling phase