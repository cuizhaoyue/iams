package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

const (
	flagHelp          = "help"
	flagHelpShorthand = "h"
)

func helpCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command.",
		Long: `Help provides help for any command in the application.
simply type ` + name + ` help [path to command] for full details.`,

		Run: func(c *cobra.Command, args []string) {
			// 查找root Command的子命令
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unkonwn help topic %#q\n", args)
				_ = c.Root().Usage()
			} else {
				cmd.InitDefaultHelpFlag() // 为子命令添加默认的help flag
				_ = cmd.Help()            // 打印help信息
			}
		},
	}
}

func addHelpCommandFlag(usage string, fs *pflag.FlagSet) {
	fs.BoolP(
		flagHelp,
		flagHelpShorthand,
		false,
		fmt.Sprintf("Help for the %s command.", color.GreenString(strings.Split(usage, " ")[0])),
	)
}
