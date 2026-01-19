package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(level, format string) {
	// 设置日志级别
	var zerologLevel zerolog.Level
	switch level {
	case "debug":
		zerologLevel = zerolog.DebugLevel
	case "info":
		zerologLevel = zerolog.InfoLevel
	case "warn":
		zerologLevel = zerolog.WarnLevel
	case "error":
		zerologLevel = zerolog.ErrorLevel
	default:
		zerologLevel = zerolog.InfoLevel
	}

	// 设置时间格式
	zerolog.TimeFieldFormat = time.RFC3339

	// 根据 format 选择输出格式
	if format == "json" {
		log = zerolog.New(os.Stdout).
			Level(zerologLevel).
			With().
			Timestamp().
			CallerWithSkipFrameCount(3).
			Logger()
	} else {
		// console 格式，更易读
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02T15:04:05.000Z07:00",
		}
		log = zerolog.New(output).
			Level(zerologLevel).
			With().
			Timestamp().
			CallerWithSkipFrameCount(3).
			Logger()
	}
}

func Info(args ...interface{})  { log.Info().Msg(fmt.Sprint(args...)) }
func Debug(args ...interface{}) { log.Debug().Msg(fmt.Sprint(args...)) }
func Warn(args ...interface{})  { log.Warn().Msg(fmt.Sprint(args...)) }
func Error(args ...interface{}) { log.Error().Msg(fmt.Sprint(args...)) }
func Fatal(args ...interface{}) { log.Fatal().Msg(fmt.Sprint(args...)) }

func Infof(template string, args ...interface{})  { log.Info().Msgf(template, args...) }
func Debugf(template string, args ...interface{}) { log.Debug().Msgf(template, args...) }
func Warnf(template string, args ...interface{})  { log.Warn().Msgf(template, args...) }
func Errorf(template string, args ...interface{}) { log.Error().Msgf(template, args...) }
func Fatalf(template string, args ...interface{}) { log.Fatal().Msgf(template, args...) }
