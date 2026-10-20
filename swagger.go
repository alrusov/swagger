package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"gopkg.org/swaggerui"

	"github.com/alrusov/config"
	"github.com/alrusov/initializer"
	"github.com/alrusov/misc"
	"github.com/alrusov/stdhttp"
)

//----------------------------------------------------------------------------------------------------------------------------//

type (
	Config struct {
		OpenAPI string `toml:"openapi"`
		h       http.Handler
	}
)

//----------------------------------------------------------------------------------------------------------------------------//

func init() {
	// Регистрируем инициализатор
	initializer.RegisterModuleInitializer(initModule)
}

// Инициализация
func initModule(appCfg any, h any) (err error) {
	cfg, err := config.GetBlock[*Config](appCfg, "Swagger")
	if err != nil {
		return
	}

	if cfg.OpenAPI == "" {
		return
	}

	hdl, ok := h.(*stdhttp.HTTP)
	if !ok {
		err = fmt.Errorf("http handler is %T, expected %T", h, hdl)
		return
	}

	bundle := swaggerui.DefaultBundle
	bundle.URL = cfg.OpenAPI
	cfg.h = http.StripPrefix("/swagger", swaggerui.Handler(bundle))
	hdl.AddHandler(cfg, false)
	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// Проверка валидности Config
func (x *Config) Check(cfg any) (err error) {
	msgs := misc.NewMessages()
	defer msgs.Free()

	return msgs.Error()
}

//----------------------------------------------------------------------------------------------------------------------------//

func (cfg *Config) Handler(id uint64, prefix string, path string, w http.ResponseWriter, r *http.Request) (processed bool) {
	if path == "/swagger" || strings.HasPrefix(path, "/swagger/") {
		processed = true
		cfg.h.ServeHTTP(w, r)
		return
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//
