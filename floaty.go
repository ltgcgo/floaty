package floaty

import (
	"fmt"
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

// Initialize!
type FloatyItem struct {
	id        string
	Length    int   `json:"length,omitempty"`
	Duration  int64 `json:"duration,omitempty"`
	lastWrite int64
	nextWrite int64
}
type FloatyModule struct {
	logger *zap.Logger
	Values map[string]*FloatyItem `json:"values,omitempty"`
}

// Parse the Caddyfile directives
func (module *FloatyModule) UnmarshalCaddyfile(
	dispenser *caddyfile.Dispenser,
) error {
	// Initialize the maps
	if module.Values == nil {
		module.Values = make(map[string]*FloatyItem)
	}
	module.Values["rootId"] = new(FloatyItem)
	// Parse the rootId parameters
	arg1 := dispenser.NextArg()
	arg2 := dispenser.NextArg()
	var length int
	var duration int64
	if arg1 && arg2 {
		lengthRaw, err := strconv.Atoi(dispenser.Val())
		if err != nil {
			return dispenser.Err("Cannot parse length into an integer")
		}
		if lengthRaw < 4 {
			lengthRaw = 4
		} else if length > 96 {
			lengthRaw = 96
		}
		length = lengthRaw
	} else {
		length = 8
	}
	if dispenser.NextArg() {
		durationObj, err := time.ParseDuration(dispenser.Val())
		if err != nil {
			return dispenser.Err("Cannot parse duration into a duration")
		}
		duration = durationObj.Milliseconds()
		if duration < 10000 {
			duration = 10000
		}
	} else {
		duration = 5400000 // 15 minutes
	}
	module.Values["rootId"].Duration = duration
	module.Values["rootId"].Length = length
	// Parse configs for additional rolling IDs
	for dispenser.NextBlock(0) {
		mapKey := dispenser.Val()
		module.Values[mapKey] = new(FloatyItem)
		var length int
		var duration int64
		if !dispenser.NextArg() {
			length = 8
			duration = 900000
		} else {
			lengthRaw, err := strconv.Atoi(dispenser.Val())
			if err != nil {
				return dispenser.Err("Cannot parse length into an integer")
			}
			if lengthRaw < 4 {
				lengthRaw = 4
			} else if length > 96 {
				lengthRaw = 96
			}
			length = lengthRaw
			if !dispenser.NextArg() {
				duration = 900000
			} else {
				durationObj, err := time.ParseDuration(dispenser.Val())
				if err != nil {
					return dispenser.Err("Cannot parse duration into a duration")
				}
				duration = durationObj.Milliseconds()
				if duration < 10000 {
					duration = 10000
				}
			}
		}
		module.Values[mapKey].Length = length
		module.Values[mapKey].Duration = duration
	}
	return nil
}
func caddyParser(
	helper httpCaddyfile.Helper,
) (
	caddyHttp.MiddlewareHandler,
	error,
) {
	module := new(FloatyModule)
	err := module.UnmarshalCaddyfile(helper.Dispenser)
	if err == nil {
		fmt.Println("\x1b[1;33m[Floaty]\x1b[0;m No errors are present in the Caddyfile.")
	}
	return module, err
}

// Register the module
func (FloatyModule) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.floaty",
		New: func() caddy.Module {
			return new(FloatyModule)
		},
	}
}
func init() {
	caddy.RegisterModule(FloatyModule{})
	httpCaddyfile.RegisterHandlerDirective("floaty", caddyParser)
}

// Provision!
func (module *FloatyModule) Provision(ctx caddy.Context) error {
	timeNow := time.Now().UnixMilli()
	if module.Values == nil {
		module.Values = make(map[string]*FloatyItem)
		module.Values["rootId"] = new(FloatyItem)
		module.Values["rootId"].Duration = 900000
		module.Values["rootId"].Length = 8
		fmt.Println("\x1b[1;33m[Floaty]\x1b[0;m Map not yet parsed before provision. Creating the map.")
	}
	for _, mapConf := range module.Values {
		mapConf.id = nanoid.Must(mapConf.Length)
		mapConf.lastWrite = timeNow
		mapConf.nextWrite = timeNow + mapConf.Duration
	}
	module.logger = ctx.Logger()
	return nil
}

// Validate!
/*func (module *FloatyModule) Validate() error {
	return nil;
}*/

// Handle!
// Handle requests with placeholder replacements
func (module FloatyModule) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
	handler caddyHttp.Handler,
) error {
	timeNow := time.Now().UnixMilli()
	repl := request.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	// Refresh IDs when stale
	// Refresh root ID
	for mapKey, mapConf := range module.Values {
		if mapConf.nextWrite <= timeNow {
			module.logger.Info(
				"A Floaty ID has expired! Current state.",
				zap.String("key", mapKey),
				zap.String("id", mapConf.id),
				zap.Int64("lastWrite", mapConf.lastWrite),
				zap.Int64("nextWrite", mapConf.nextWrite),
			)
			mapConf.id = nanoid.Must(mapConf.Length)
			mapConf.lastWrite = timeNow
			mapConf.nextWrite = timeNow + mapConf.Duration
			module.logger.Info(
				"A Floaty ID has expired! New state.",
				zap.String("key", mapKey),
				zap.String("id", mapConf.id),
				zap.Int64("lastWrite", mapConf.lastWrite),
				zap.Int64("nextWrite", mapConf.nextWrite),
			)
		}
	}
	// Set values for placeholders
	repl.Set("http.floaty", module.Values["rootId"].id)
	for mapKey, mapConf := range module.Values {
		if mapKey == "rootId" {
			continue
		}
		repl.Set("http.floaty." + mapKey, mapConf.id)
	}
	return handler.ServeHTTP(writer, request)
}

// Interface guards
var (
	_ caddy.Provisioner           = (*FloatyModule)(nil)
	_ caddyfile.Unmarshaler       = (*FloatyModule)(nil)
	_ caddyHttp.MiddlewareHandler = (*FloatyModule)(nil)
)
