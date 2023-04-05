package app

import (
	"fmt"
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/fatih/color"
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/marmotedu/component-base/pkg/cli/globalflag"
	"github.com/marmotedu/component-base/pkg/term"
	"github.com/marmotedu/component-base/pkg/version"
	"github.com/marmotedu/component-base/pkg/version/verflag"
	"github.com/marmotedu/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	progressMessage = color.GreenString("===>")
)

// App 是客户端应用的主要结构体.
// 建议使用app.NewApp()函数创建app.
type App struct {
	basename    string               // 应用的二进制文件名
	name        string               // 应用的简短描述
	description string               // 应用的详细描述
	options     CliOptions           // 应用的命令行选项
	runFunc     RunFunc              // 应用的启动函数，初始化应用，并最终启动 HTTP 和 GRPC Web 服务
	silence     bool                 // 是否启动静默模式
	noVersion   bool                 // 是否提供version flag
	noConfig    bool                 // 是否提供config flag
	commands    []*Command           // 应用的子命令
	args        cobra.PositionalArgs // 定义应用的位置参数校验函数
	cmd         *cobra.Command       // 应用的root Command
}

// Option 定义了初始化应用结构体时的可选参数.
type Option func(*App)

// WithOptions 用于打开应用程序的函数以从命令行读取或从配置文件读取参数。
func WithOptions(opt CliOptions) Option {
	return func(app *App) {
		app.options = opt
	}
}

// RunFunc 定义了应用的启动调用函数.
type RunFunc func(basename string) error

// WithRunFunc 用于设置应用启动的回调函数.
func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.runFunc = run
	}
}

// WithDescription 用于设置应用的描述信息.
func WithDescription(desc string) Option {
	return func(app *App) {
		app.description = desc
	}
}

// WithSilence  将应用程序设置为静默模式.
// 在该模式下，程序启动信息、配置信息和版本信息不会打印在控制台中.
func WithSilence() Option {
	return func(app *App) {
		app.silence = true
	}
}

// WithNoVersion 设置应用不提供version flag.
func WithNoVersion() Option {
	return func(app *App) {
		app.noVersion = true
	}
}

// WithNoConfig 设置应用不提供config flag.
func WithNoConfig() Option {
	return func(app *App) {
		app.noConfig = true
	}
}

// WithValidArgs 设置校验函数以验证非flag参数.
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(app *App) {
		app.args = args
	}
}

// WithDefaultValidArgs 设置校验non-flag参数的默认验证函数.
func WithDefaultValidArgs() Option {
	return func(app *App) {
		app.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

func NewApp(name string, basename string, opts ...Option) *App {
	app := &App{
		name:     name,
		basename: basename,
	}

	// 设置可选参数
	for _, opt := range opts {
		opt(app)
	}

	app.buildCommand()

	return app
}

// 构建应用的root Command
func (a *App) buildCommand() {
	cmd := &cobra.Command{
		Use:           a.basename,
		Short:         a.name,
		Long:          a.description,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	cmd.Flags().SortFlags = true // 打印usage信息时对flag排序
	InitFlags(cmd.Flags())

	// 添加子命令
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(a.basename))
	}

	// 获取为此应用设置的所有flags，并添加到当前应用的FlagSet中
	var namedFlagSets cliflag.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	// 把version flag划分到global组
	if !a.noVersion {
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}

	// 添加config flag，划分到global分组，并添加读取配置文件的初始化操作
	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}

	// 为global flag分组添加help flag
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())

	// 将global FlagSet添加到此应用的FlagSet中
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	// 设置执行应用启用函数的运行函数
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	// 设置cmd的usage/help函数
	addCmdTemplate(cmd, namedFlagSets)

	a.cmd = cmd
}

// Run 用于启动应用程序.
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// Command 返回应用中的*cobra.Command实例
func (a *App) Command() *cobra.Command {
	return a.cmd
}

// 设置执行*cobra.Command.RunE的运行函数
func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir() // 打印工作目录

	// 打印应用的所有flag
	PrintFlags(cmd.Flags())

	// 添加version flag，打印version信息.
	if !a.noVersion {
		verflag.PrintAndExitIfRequested()
	}

	// 添加配置
	if !a.noConfig {
		// 把从命令行传入的所有flag的值绑定到读取的配置中
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		// 把读取到的配置的值反序列化为Options结构体的值
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	// 非静默模式时打印版本和配置信息
	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, a.name)
		if !a.noVersion {
			log.Infof("%v Version: `%s`", progressMessage, version.Get().ToJSON())
		}
		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}

	// 执行options的补全和验证操作
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

// 执行options的补全操作和验证操作，并打印配置信息
func (a *App) applyOptionRules() error {
	if completeableOptions, ok := a.options.(CompletableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	if errs := a.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}

	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	}

	return nil
}

// 打印当前工作目录
func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}

// 设置cmd的help/usage
func addCmdTemplate(cmd *cobra.Command, namedFlagSets cliflag.NamedFlagSets) {
	usageFmt := "Usage:\n %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())

	// 设置把Usage信息写入标准输出或错误输出
	cmd.SetUsageFunc(func(command *cobra.Command) error {
		fmt.Fprintf(command.OutOrStdout(), usageFmt, command.UseLine())
		cliflag.PrintSections(command.OutOrStderr(), namedFlagSets, cols)

		return nil
	})

	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		fmt.Fprintf(command.OutOrStdout(), "%s\n\n"+usageFmt, command.Long, command.UseLine())
		cliflag.PrintSections(command.OutOrStderr(), namedFlagSets, cols)
	})
}
