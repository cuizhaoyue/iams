package app

import (
	"fmt"
	"github.com/marmotedu/component-base/pkg/util/homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const configFlagName = "config"

var cfgFile string

// nolint: gochecknoinits
func init() {
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, TOML, YAML, HCL or Java properties formats.")
}

// 添加config flag到指定的FlagSet中，读取环境变量和配置文件
func addConfigFlag(basename string, fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(configFlagName))

	// 读取环境变量
	viper.AutomaticEnv()
	// 设置环境变量前缀
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(basename), "-", "_", -1))
	// 设置环境变量的
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// 设置执行*cobra.Command的Execute方法之前要运行的函数
	cobra.OnInitialize(func() {
		// 读取配置文件
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			// 设置配置文件路径
			viper.AddConfigPath(".")

			if names := strings.Split(basename, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join(homedir.HomeDir(), "."+names[0]))
				viper.AddConfigPath(filepath.Join("/etc", names[0]))
			}

			// 设置配置文件的名称
			viper.SetConfigName(basename)
		}

		// 加载配置文件
		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", cfgFile, err)
			os.Exit(1)
		}
	})
}
