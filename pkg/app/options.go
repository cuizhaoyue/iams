package app

import (
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
)

// CliOptions 提取了用于从命令行读取参数的配置选项。
type CliOptions interface {
	Flags() cliflag.NamedFlagSets
	Validate() []error
}

// CompletableOptions 提取了可以被执行completed的选项
type CompletableOptions interface {
	Complete() error
}

// PrintableOptions 提取可以被打印的options
type PrintableOptions interface {
	String() string
}
