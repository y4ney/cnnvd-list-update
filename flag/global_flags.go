package flag

import (
	"github.com/0yaney0/cnnvd-list-update/utils"
	"github.com/spf13/cobra"
	"time"
)

var (
	ConfigFileFlag = Flag{
		Name:       "config",
		ConfigName: "config",
		Shorthand:  "c",
		Value:      "config.yaml",
		Usage:      "指定config配置文件的路径",
		Persistent: true,
	}
	ShowVersionFlag = Flag{
		Name:       "version",
		ConfigName: "version",
		Shorthand:  "v",
		Value:      false,
		Usage:      "展示版本信息",
		Persistent: true,
	}
	QuietFlag = Flag{
		Name:       "quiet",
		ConfigName: "quiet",
		Shorthand:  "q",
		Value:      false,
		Usage:      "关闭进度条和日志输出",
		Persistent: true,
	}
	DebugFlag = Flag{
		Name:       "debug",
		ConfigName: "debug",
		Shorthand:  "d",
		Value:      false,
		Usage:      "指定 debug 模式",
		Persistent: true,
	}
	TimeoutFlag = Flag{
		Name:       "timeout",
		ConfigName: "timeout",
		Value:      time.Second * 300, // 5 mins
		Usage:      "指定超时时间",
		Persistent: true,
	}
	CacheDirFlag = Flag{
		Name:       "cache-dir",
		ConfigName: "cache.dir",
		Value:      utils.DefaultCacheDir(),
		Usage:      "指定缓存目录",
		Persistent: true,
	}
)

// GlobalFlagGroup composes global flags
type GlobalFlagGroup struct {
	ConfigFile  *Flag
	ShowVersion *Flag // spf13/cobra can't override the logic of version printing like VersionPrinter in urfave/cli. -v needs to be defined ourselves.
	Quiet       *Flag
	Debug       *Flag
	Timeout     *Flag
	CacheDir    *Flag
}

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type GlobalOptions struct {
	ConfigFile  string
	ShowVersion bool
	Quiet       bool
	Debug       bool
	Timeout     time.Duration
	CacheDir    string
}

func NewGlobalFlagGroup() *GlobalFlagGroup {
	return &GlobalFlagGroup{
		ConfigFile:  &ConfigFileFlag,
		ShowVersion: &ShowVersionFlag,
		Quiet:       &QuietFlag,
		Debug:       &DebugFlag,
		Timeout:     &TimeoutFlag,
		CacheDir:    &CacheDirFlag,
	}
}
func (f *GlobalFlagGroup) flags() []*Flag {
	return []*Flag{f.ConfigFile, f.ShowVersion, f.Quiet, f.Debug, f.Timeout, f.CacheDir}
}
func (f *GlobalFlagGroup) AddFlags(cmd *cobra.Command) {
	for _, flag := range f.flags() {
		addFlag(cmd, flag)
	}
}
func (f *GlobalFlagGroup) Bind(cmd *cobra.Command) error {
	for _, flag := range f.flags() {
		if err := bind(cmd, flag); err != nil {
			return err
		}
	}
	return nil
}
func (f *GlobalFlagGroup) ToOptions() GlobalOptions {
	return GlobalOptions{
		ConfigFile:  getString(f.ConfigFile),
		ShowVersion: getBool(f.ShowVersion),
		Quiet:       getBool(f.Quiet),
		Debug:       getBool(f.Debug),
		Timeout:     getDuration(f.Timeout),
		CacheDir:    getString(f.CacheDir),
	}
}
