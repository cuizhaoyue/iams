package log

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	flagName              = "log.name"
	flagLevel             = "log.level"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"

	consoleFormat = "console"
	jsonFormat    = "json"
)

// Options contains configuration items related to log.
type Options struct {
	Name              string   `json:"name" mapstructure:"name"`                             // Logger 的名字。
	Level             string   `json:"level" mapstructure:"level"`                           // 日志级别，优先级从低到高依次为：Debug, Info, Warn, Error, Dpanic, Panic, Fatal。
	Format            string   `json:"format" mapstructure:"format"`                         // 支持的日志输出格式，目前支持 Console 和 JSON 两种。Console 其实就是 Text 格式。
	Development       bool     `json:"development" mapstructure:"development"`               // 是否是开发模式。如果是开发模式，会对 DPanicLevel 进行堆栈跟踪。
	EnableColor       bool     `json:"enable-color" mapstructure:"enable-color"`             // 是否开启颜色输出，true，是；false，否。
	DisableCaller     bool     `json:"disable-caller" mapstructure:"disable-caller"`         // 是否开启 caller，如果开启会在日志中显示调用日志所在的文件、函数和行号。
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"` // 是否在 Panic 及以上级别禁止打印堆栈信息。
	OutputPaths       []string `json:"output-paths" mapstructure:"output-paths"`             // 支持输出到多个输出，用逗号分开。支持输出到标准输出（stdout）和文件。
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"` // zap 内部 (非业务) 错误日志输出路径，多个输出，用逗号分开。
}

func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		Format:            consoleFormat,
		Development:       false,
		EnableColor:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
}

// Validate validate the options fields.
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}

// Build 根据Config和Optons构建一个全局的zap logger.
func (o *Options) Build() error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	levelEncoder := zapcore.CapitalLevelEncoder
	if o.Format == consoleFormat && o.EnableColor {
		levelEncoder = zapcore.CapitalColorLevelEncoder
	}

	zc := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: o.Format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    levelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: milliSecondDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      o.OutputPaths,
		ErrorOutputPaths: o.ErrorOutputPaths,
	}

	logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		return err
	}
	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}
