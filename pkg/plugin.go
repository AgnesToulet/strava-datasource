package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/grafana/grafana-plugin-model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	cache "github.com/patrickmn/go-cache"
)

const DATA_PATH_VARIABLE = "GF_STRAVA_DS_DATA_PATH"

var pluginLogger = hclog.New(&hclog.LoggerOptions{
	Name:  "strava-datasource",
	Level: hclog.LevelFromString("DEBUG"),
})

func main() {
	pluginLogger.Debug("Running Strava backend datasource")

	var dataDir string
	dataDir, exist := os.LookupEnv(DATA_PATH_VARIABLE)
	if !exist {
		dataDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	pluginLogger.Debug("Plugin data dir", "path", dataDir)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "grafana_plugin_type",
			MagicCookieValue: "datasource",
		},
		Plugins: map[string]plugin.Plugin{
			"strava-backend-datasource": &datasource.DatasourcePluginImpl{Plugin: &StravaPlugin{
				datasourceCache: cache.New(10*time.Minute, 10*time.Minute),
				logger:          pluginLogger,
				dataDir:         dataDir,
			}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
