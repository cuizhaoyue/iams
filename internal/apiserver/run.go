package apiserver

import "github.com/cuizhaoyue/iams/internal/apiserver/config"

// Run 运行指定的APIServer.
func Run(cfg *config.Config) error {
	server, err := createAPIServer(cfg)

	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
