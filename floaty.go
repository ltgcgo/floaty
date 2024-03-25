package floaty

import (
	"net/http"
	//"strconv"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	httpCaddyfile "github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	caddyHttp "github.com/caddyserver/caddy/v2/modules/caddyhttp"
	nanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

// Initialize!
type FloatyItem struct {
	id string
	length int `json:"length"`
	duration int64 `json:"duration"`
	lastWrite int64
	nextWrite int64
}
type FloatyModule struct {
	logger *zap.Logger
	values map[string]*FloatyItem
}
// Parse the Caddyfile directives
func (module *FloatyModule) UnmarshalCaddyfile(
	dispenser *caddyfile.Dispenser,
) error {
	return nil;
}
func caddyParser(
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
// Register the module
func (FloatyModule) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.floaty",
		New: func() caddy.Module {
			return new(FloatyModule)
		},
	};
}
func init() {
	caddy.RegisterModule(FloatyModule{});
	httpCaddyfile.RegisterHandlerDirective("floaty", caddyParser);
}

// Provision!
func (module *FloatyModule) Provision(ctx caddy.Context) error {
	timeNow := time.Now().UnixMilli();
	module.logger = ctx.Logger();
	module.values = make(map[string]*FloatyItem);
	module.values["rootId"] = new(FloatyItem);
	module.values["rootId"].duration = 5000;
	module.values["rootId"].length = 16;
	module.values["rootId"].id = nanoid.Must(module.values["rootId"].length);
	module.values["rootId"].lastWrite = timeNow;
	module.values["rootId"].nextWrite = timeNow + module.values["rootId"].duration;
	module.logger.Info(
		"Floaty has been provisioned!",
		zap.String("rootId", module.values["rootId"].id),
		zap.Int64("lastWrite", module.values["rootId"].lastWrite),
		zap.Int64("nextWrite", module.values["rootId"].nextWrite),
	);
	return nil;
}

// Validate!
func (module *FloatyModule) Validate() error {
	return nil;
}

// Handle!
// Handle requests with placeholder replacements
func (module FloatyModule) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
	handler caddyHttp.Handler,
) error {
	timeNow := time.Now().UnixMilli();
	repl := request.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer);
	// Refresh IDs when stale
	// Refresh root ID
	if (module.values["rootId"].nextWrite <= timeNow) {
		module.logger.Info(
			"Root ID of Floaty has expired! Current state.",
			zap.String("id", module.values["rootId"].id),
			zap.Int64("lastWrite", module.values["rootId"].lastWrite),
			zap.Int64("nextWrite", module.values["rootId"].nextWrite),
			zap.Int64("timeNow", timeNow),
		);
		module.values["rootId"].id = nanoid.Must(module.values["rootId"].length);
		module.values["rootId"].lastWrite = timeNow;
		module.values["rootId"].nextWrite = timeNow + module.values["rootId"].duration;
		module.logger.Info(
			"Root ID of Floaty has expired! New state.",
			zap.String("id", module.values["rootId"].id),
			zap.Int64("lastWrite", module.values["rootId"].lastWrite),
			zap.Int64("nextWrite", module.values["rootId"].nextWrite),
		);
	};
	// Set values for placeholders
	repl.Set("http.floaty", module.values["rootId"].id);
	/*module.logger.Info(
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
