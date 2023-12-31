// This plugin used luludotdev/caddy-requestid as a starting point

package floaty

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	httpCaddyfile "github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	caddyHttp "github.com/caddyserver/caddy/v2/modules/caddyhttp"
	nanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

// Initialization phase

// Floaty implements global placeholders that rolls with a set interval
type FloatyPool struct {
	lastProvision int64
	nextProvision int64
	map map[string]string
}
type FloatyModule struct {
	// Logger
	logger *zap.Logger

	// generated service name
	Service string

	// Length of instance ID
	Length int `json:"length"`

	// Duration of validity
	Duration time.Duration `json:"duration"`

	// Map of additional instance IDs to be set
	Additional map[string]int `json:"additional,omitempty"`

	// Initialized master ID
	InstanceId map[string]string

	// Map of additional instance IDs initialized
	MappedIds map[string]FloatyPool
}

// Caddyfile syntax parsing
func (module *FloatyModule) UnmarshalCaddyfile(
	dispenser *caddyfile.Dispenser,
) error {
	arg1 := dispenser.NextArg();
	arg2 := dispenser.NextArg();
	// Standalone length parsing
	if arg1 && arg2 {
		value := dispenser.Val();
		length, err := strconv.Atoi(value);
		if err != nil {
			return dispenser.Err("Conversion of length to integer failed.");
		};
		if length < 1 {
			return dispenser.Err("Length must be a positive integer.");
		};
		module.Length = length;
	};
	// Mapped IDs length parsing
	return nil;
}

// Entrypoint for Caddyfile parsing
func cfParser(
	helper httpCaddyfile.Helper,
) (
	caddyHttp.MiddlewareHandler,
	error,
) {
	module := new(FloatyModule);
	err := module.UnmarshalCaddyfile(helper.Dispenser);
	if (err != nil) {
		return nil, err;
	};
	return module, nil;
}

// Initialize the module
func init() {
	caddy.RegisterModule(FloatyModule{});
	httpCaddyfile.RegisterHandlerDirective("floaty", cfParser);
}

// Register the Caddy plugin
func (FloatyModule) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.floaty",
		New: func() caddy.Module {
			return new(FloatyModule)
		},
	};
}

// Provisioning phase

// Set up the IDs
var FloatyModuleGlobal string;
var FloatyModuleMapped map[string]string;
func (module *FloatyModule) Provision(ctx caddy.Context) error {
	// Create a logger
	module.logger = ctx.Logger();
	// Normalize the parameters
	if module.Length < 1 {
		module.Length = 12;
	};
	if module.Additional == nil {
		module.Additional = make(map[string]int);
	};
	// Generate the IDs
	FloatyModuleGlobal = nanoid.Must(module.Length);
	for i0, e0 := range module.Additional {
		if e0 < 1 {
			e0 = 12;
		};
		FloatyModuleMapped[i0] = nanoid.Must(e0);
	};
	// Bind the IDs to global variables
	module.InstanceId = FloatyModuleGlobal;
	module.MappedIds = FloatyModuleMapped;
	// Log the created IDs
	module.logger.Info(
		"Floaty has provisioned global ID: ",
		zap.String("id", module.InstanceId),
	);
	return nil;
}

// Handling phase

// Handle requests with placeholder replacements
func (module FloatyModule) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
	handler caddyHttp.Handler,
) error {
	repl := request.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer);
	// Set values for placeholders
	/*repl.Set("http.floaty", module.InstanceId);
	module.logger.Info(
		"Floaty has accessed global ID: ",
		zap.String("id", module.InstanceId),
	);
	for i0, e0 := range module.MappedIds {
		repl.Set("http.floaty." + i0, e0);
	};*/
	return handler.ServeHTTP(writer, request);
}

// Interface guards
var (
	_ caddy.Provisioner = (*FloatyModule)(nil)
	_ caddyfile.Unmarshaler = (*FloatyModule)(nil)
	_ caddyHttp.MiddlewareHandler = (*FloatyModule)(nil)
);
