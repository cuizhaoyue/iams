package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

// Command 是一个命令行应用的子命令结构体.
// 建议使用app.NewCommand()函数创建子命令.
type Command struct {
	usage    string         // 子命令的使用信息，也是子命令的名称
	desc     string         // 子命令的简短描述
	options  CliOptions     // 子命令的Flag选项
	commands []*Command     // 下层的子命令
	runFunc  RunCommandFunc // 启动时调用的函数
}

// CommandOption 定义了初始化Command结构体的可选参数
type CommandOption func(*Command)

// WithCommandOption 用于打开应用程序的函数以从命令行读取参数。
func WithCommandOption(opts CliOptions) CommandOption {
	return func(command *Command) {
		command.options = opts
	}
}

// RunCommandFunc 定义了应用的命令启动时的调用函数.
type RunCommandFunc func(args []string) error

// WithCommandRunFunc 用来设置应用的命令启动时的调用函数.
func WithCommandRunFunc(run RunCommandFunc) CommandOption {
	return func(command *Command) {
		command.runFunc = run
	}
}

// NewCommand 基于给定的命令名称和其它选项创建一个新的子命令实例.
func NewCommand(usage string, desc string, opts ...CommandOption) *Command {
	c := &Command{
		usage: usage,
		desc:  desc,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// AddCommand 添加子命令到当前命令
func (c *Command) AddCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

// AddCommands 添加多个子命令到当前命令.
func (c *Command) AddCommands(cmds ...*Command) {
	c.commands = append(c.commands, cmds...)
}

// 根据自定义的Command结构体构建*cobra.Command实例
func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}
	cmd.SetOut(os.Stdout)        // 设置标准输出的目标
	cmd.Flags().SortFlags = true // 对子命令的flag排序

	// 添加当前子命令的下层命令
	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}

	// 设置Command的启动函数
	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}

	// 添加子命令的flag选项
	if c.options != nil {
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
	}

	// 为子命令添加help flag
	addHelpCommandFlag(c.usage, cmd.Flags())

	return cmd
}

// 构建*cobra.Command的Run函数，运行该函数出错时退出程序.
func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

// AddCommand 添加子命令到应用.
func (a *App) AddCommand(cmd *Command) {
	a.commands = append(a.commands, cmd)
}

// AddCommands 添加多个子命令到应用.
func (a *App) AddCommands(cmds ...*Command) {
	a.commands = append(a.commands, cmds...)
}
