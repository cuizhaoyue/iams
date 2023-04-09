package apiserver

import "github.com/cuizhaoyue/iams/pkg/shutdown"

// apiserver应用配置
type apiServer struct {
	gs *shutdown.GracefuleShutdown // 负责服务优雅关闭
	// redisOptions *

}
