package golog

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// GlobalLevel global level
var GlobalLevel = DEBUG

const (
	WriterNameConsole = "console_writer"
	WriterNameFile    = "file_writer"
)

// LogConfig log config
type LogConfig struct {
	//Global level
	Level string `json:"level" mapstructure:"level"`
	Debug bool   `json:"debug" mapstructure:"debug"` // output log info or not for go-log
	//If display full path of the file which log belongs to
	FullPath      bool                 `json:"full_path" mapstructure:"full_path"`
	ConsoleWriter ConsoleWriterOptions `json:"console_writer" mapstructure:"console_writer"`
	FileWriter    FileWriterOptions    `json:"file_writer" mapstructure:"file_writer"`
}

// SetupLog setup log
func SetupLog(lc LogConfig) (err error) {
	if !lc.Debug {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stdout)
	}

	// global config
	GlobalLevel = getLevel(lc.Level)

	// writer enable
	// 1. if not set level, use global level;
	// 2. if set level, use min level
	validGlobalMinLevel := ACCESS // default max level
	validGlobalMinLevelBy := "global"

	fileWriterLevelDefault := GlobalLevel
	consoleWriterLevelDefault := GlobalLevel

	if lc.ConsoleWriter.Enable {
		consoleWriterLevelDefault = getLevelDefault(lc.ConsoleWriter.Level, GlobalLevel, WriterNameConsole)
		validGlobalMinLevel = maxInt(consoleWriterLevelDefault, validGlobalMinLevel)
		if validGlobalMinLevel == consoleWriterLevelDefault {
			validGlobalMinLevelBy = WriterNameConsole
		}
	}

	if lc.FileWriter.Enable {
		fileWriterLevelDefault = getLevelDefault(lc.FileWriter.Level, GlobalLevel, WriterNameFile)
		validGlobalMinLevel = maxInt(fileWriterLevelDefault, validGlobalMinLevel)
		if validGlobalMinLevel == fileWriterLevelDefault {
			validGlobalMinLevelBy = WriterNameFile
		}
	}

	fullPath := lc.FullPath
	WithFullPath(fullPath)
	SetLevel(validGlobalMinLevel)

	if lc.ConsoleWriter.Enable {
		w := NewConsoleWriterWithOptions(lc.ConsoleWriter)
		w.level = consoleWriterLevelDefault
		log.Printf("[go-log] enable " + WriterNameConsole + " with level " + LevelFlags[consoleWriterLevelDefault])
		Register(w)
	}

	if lc.FileWriter.Enable {
		w := NewFileWriterWithOptions(lc.FileWriter)
		w.level = fileWriterLevelDefault
		log.Printf("[go-log] enable    " + WriterNameFile + " with level " + LevelFlags[fileWriterLevelDefault])
		Register(w)
	}

	log.Printf("[go-log] valid global_level(min:%v, flag:%v, by:%v), default(%v, flag:%v)",
		validGlobalMinLevel, LevelFlags[validGlobalMinLevel], validGlobalMinLevelBy, GlobalLevel, LevelFlags[GlobalLevel])
	return nil
}

// SetLogWithConf setup log with config file
func SetLogWithConf(file string) (err error) {
	var lc LogConfig
	cnt, err := ioutil.ReadFile(file)

	if err = json.Unmarshal(cnt, &lc); err != nil {
		return
	}
	return SetupLog(lc)
}

// SetLog setup log with config []byte
func SetLog(config []byte) (err error) {
	var lc LogConfig
	if err = json.Unmarshal(config, &lc); err != nil {
		return
	}
	return SetupLog(lc)
}

func getLevel(flag string) int {
	return getLevelDefault(flag, DEBUG, "")
}

// maxInt return max int
func maxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}
