package app

import (
	goflag "flag"
	"github.com/marmotedu/log"
	"github.com/spf13/pflag"
	"strings"
)

// InitFlags 规范化、解析然后记录命令行标志
func InitFlags(flags *pflag.FlagSet) {
	// 规范化flag名称
	flags.SetNormalizeFunc(WordSepNormalizeFunc)
	// 把标准库中的flag包设置的FlagSet添加到pflag.FlagSet中
	flags.AddGoFlagSet(goflag.CommandLine)
}

// WordSepNormalizeFunc 更改所有包含"_"分隔符的flag
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}

	return pflag.NormalizedName(name)
}

func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		log.Debugf("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
