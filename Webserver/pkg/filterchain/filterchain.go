package filterchain

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"webserver/Webserver/pkg/http_manager"
	"webserver/Webserver/pkg/listener"
	"webserver/Webserver/pkg/tls"
	"webserver/Webserver/pkg/types"
)

func BuildFilterChain(configPath string) (types.Filter, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var rawConfig struct {
		FilterChain map[string]map[string]interface{} `yaml:"filterChain"`
	}

	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, err
	}

	var filters []types.Filter
	var first, prev types.Filter

	for name, cfg := range rawConfig.FilterChain {
		var filter types.Filter
		switch name {
		case "listener":
			filter = &listener.ListenerFilter{}
		case "tls_socket":
			// Key Provider
			certPath := cfg["certPath"].(string)
			keyPath := cfg["keyPath"].(string)
			keyProvider := tls.NewKeyProvider(certPath, keyPath)

			tlsCfg := tls.NewTLSConfigBuilder(keyProvider)
			builtCfg, err := tlsCfg.BuildDefaultTLSConfig()
			if err != nil {
				return nil, err
			}
			tlsFilter := &tls.TLSSocketFilter{}
			tlsFilter.SetTLSConfig(builtCfg)
			filter = tlsFilter
		case "httpManager":
			filter = &http_manager.HTTPManager{}
			err := filter.Init(cfg)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown filter type: %s", name)
		}

		if err := filter.Init(cfg); err != nil {
			return nil, err
		}

		if prev != nil {
			prev.SetNext(filter)
		} else {
			first = filter
		}
		prev = filter
	}

	return first, nil
}
