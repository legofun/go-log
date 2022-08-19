package golog

import (
	"testing"
	"time"
)

var (
	logConfig = `{
  "level": "common",
  "full_path": true,
  "debug": true,
	
  "file_writer": {
    "level": "abnormal",
    "filename": "./test/golog-test-%Y%M%D.log",
	"enable": true
  },

  "console_writer": {
    "level": "error",
    "enable": true,
    "color": true,
	"full_color": true
  }
}
`
)

func TestConfig(t *testing.T) {
	if err := SetLog([]byte(logConfig)); err != nil {
		panic(err)
	}
	var name = "go-log config test"
	Debug("go-log by %s debug", name)
	Common("go-log by %s common", name)
	Abnormal("go-log by %s abnormal", name)
	Transaction("go-log by %s transaction", name)
	Error("go-log by %s error", name)
	Access("go-log by %s access", name)

	time.Sleep(1 * time.Second)
}
