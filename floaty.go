// This plugin used luludotdev/caddy-requestid as a starting point

package floaty

import (
	"net/http"
	//"strconv"

	"github.com/caddyserver/caddy/v2"
	//"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	//"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	caddyHttp "github.com/caddyserver/caddy/v2/modules/caddyhttp"
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
		ID: "http.handlers.ltgc.floaty",
		New: func() caddy.Module {
			return new(FloatyID)
		},
	};
}

// Provisioning phase

// Set up the IDs
var floatyIdGlobal string;
var floatyIdMapped map[string]string;
func (module *FloatyID) Provision(ctx caddy.Context) error {
	// Normalize the parameters
	if module.Length < 1 {
		module.Length = 8;
	};
	if module.Additional == nil {
		module.Additional = make(map[string]int);
	};
	// Generate the IDs
	floatyIdGlobal = nanoid.Must(module.Length);
	for i0, e0 := range module.Additional {
		if e0 < 1 {
			e0 = 8;
		};
		floatyIdMapped[i0] = nanoid.Must(e0);
	};
	// Bind the IDs to global variables
	module.InstanceId = floatyIdGlobal;
	module.MappedIds = floatyIdMapped;
	return nil;
}

// Handling phase

// Handle requests with placeholder replacements
func (module FloatyID) ServeHTTP (
	writer http.ResponseWriter,
	request *http.Request,
	handler caddyHttp.Handler,
) error {
	repl := request.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer);
	// Set values for placeholders
	repl.Set("http.floaty", module.InstanceId);
	for i0, e0 := range module.MappedIds {
		repl.Set("http.floaty." + i0, e0);
	};
	return handler.ServeHTTP(writer, request);
}

// Caddyfile handling
