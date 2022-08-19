package golog

import (
	"testing"
)

func generateNewConsoleWriterWithOptions(level string, color, fullColor bool) *ConsoleWriter {
	options := ConsoleWriterOptions{
		Level:     level,
		Color:     color,
		FullColor: fullColor,
	}
	w := NewConsoleWriterWithOptions(options)
	w.SetColor(color)
	w.SetFullColor(fullColor)
	return w
}

func generateRegisterConsoleWriter(lg *Logger, w *ConsoleWriter, fullPath, funcName bool, layout string) {
	lg.Register(w)
	if layout == "" {
		lg.SetLayout("2006-01-02 15:04:05")
	} else {
		lg.SetLayout(layout)
	}
	lg.WithFullPath(fullPath)
	lg.WithFuncName(funcName)
}

func Test_NewConsoleWriterWithStruct(t *testing.T) {
	c := &ConsoleWriter{}
	t.Logf("%#v", c)
}

func Test_NewConsoleWriter(t *testing.T) {
	NewConsoleWriter()
}

func Test_NewConsoleWriterWithNilLogger(t *testing.T) {
	var color, fullColor, fullPath, funcName bool
	var layout string
	records := make(chan *Record, uint(0))
	close(records)
	loggerDefaultTest := newLoggerWithRecords(records)
	c := generateNewConsoleWriterWithOptions(LevelFlagDebug, color, fullColor)
	var name = "console nil logger"
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("error occur: %v", err)
			loggerDefaultTest = newLoggerWithRecords(records)
			generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
			defer loggerDefaultTest.Close()
			loggerDefaultTest.Debug("go-log by %s", name)
			loggerDefaultTest.Common("go-log by %s", name)
			loggerDefaultTest.Error("%#v", loggerDefaultTest)
		}
	}()
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)

}

func Test_NewConsoleWriterWithGlobalSet(t *testing.T) {
	var color, fullColor bool
	var layout string
	loggerDefault = NewLogger()

	defer Close()
	layout = "20060102 150405"
	SetLayout(layout)
	SetLevel(COMMON)
	WithFullPath(true)
	WithFuncName(true)
	c := generateNewConsoleWriterWithOptions(LevelFlagCommon, color, fullColor)
	Register(c)

	var name = "console with default global"
	Debug("go-log by %s", name)
	Common("go-log by %s", name)
	Abnormal("go-log by %s", name)
	Transaction("go-log by %s", name)
	Error("go-log by %s", name)
	Access("go-log by %s", name)
}

func Test_NewConsoleWriterWithLevel(t *testing.T) {
	var color, fullColor, fullPath, funcName bool
	var layout string

	records := make(chan *Record, uint(2048))
	loggerDefaultTest := newLoggerWithRecords(records)
	loggerDefaultTest.SetLevel(DEBUG)
	defer loggerDefaultTest.Close()

	c := generateNewConsoleWriterWithOptions(LevelFlagCommon, color, fullColor)
	var name = "console level"
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
	loggerDefaultTest.Debug("go-log by %s", name)
	loggerDefaultTest.Common("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by %s", name)
	loggerDefaultTest.Transaction("go-log by %s", name)
	loggerDefaultTest.Error("go-log by %s", name)
	loggerDefaultTest.Access("go-log by %s", name)
	loggerDefaultTest.Access("%#v", loggerDefaultTest)
}

func Test_NewConsoleWriterWithLevel2(t *testing.T) {
	var color, fullColor, fullPath, funcName bool
	var layout string

	records := make(chan *Record, uint(2048))
	loggerDefaultTest := newLoggerWithRecords(records)
	loggerDefaultTest.SetLevel(ABNORMAL)
	defer loggerDefaultTest.Close()

	c := generateNewConsoleWriterWithOptions(LevelFlagDebug, color, fullColor)
	var name = "console level2"
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
	loggerDefaultTest.Debug("go-log by %s", name)
	loggerDefaultTest.Common("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by fmt ", 123, " super ", name)
	loggerDefaultTest.Transaction("go-log by %s", name)
	loggerDefaultTest.Error("go-log by %s", name)
	loggerDefaultTest.Access("go-log by %s", name)
	loggerDefaultTest.Access("%#v", loggerDefaultTest)
}

func Test_NewConsoleWriterWithColor(t *testing.T) {
	var color, fullColor, fullPath, funcName bool
	var layout string

	records := make(chan *Record, uint(2048))
	loggerDefaultTest := newLoggerWithRecords(records)
	defer loggerDefaultTest.Close()

	color = true
	c := generateNewConsoleWriterWithOptions(LevelFlagDebug, color, fullColor)
	var name = "console color"
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
	loggerDefaultTest.Debug("go-log by %s", name)
	loggerDefaultTest.Common("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by %s", name)
	loggerDefaultTest.Transaction("go-log by %s", name)
	loggerDefaultTest.Error("go-log by %s", name)
	loggerDefaultTest.Access("go-log by %s", name)
	loggerDefaultTest.Access("%#v", loggerDefaultTest)
}

func Test_NewConsoleWriterWithFullColor(t *testing.T) {
	var color, fullColor, fullPath, funcName bool
	var layout string

	records := make(chan *Record, uint(2048))
	loggerDefaultTest := newLoggerWithRecords(records)
	defer loggerDefaultTest.Close()

	color = true
	fullColor = true
	c := generateNewConsoleWriterWithOptions(LevelFlagDebug, color, fullColor)
	// c := generateNewConsoleWriterWithOptions(LevelFlagEmergency, color, fullColor)
	var name = "console full color"
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
	loggerDefaultTest.Debug("go-log by %s", name)
	loggerDefaultTest.Common("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by %s", name)
	loggerDefaultTest.Transaction("go-log by %s", name)
	loggerDefaultTest.Error("go-log by %s", name)
	loggerDefaultTest.Access("go-log by %s", name)
	loggerDefaultTest.Access("%#v", loggerDefaultTest)
}

func Test_NewConsoleWriterWithFullPath(t *testing.T) {
	var color, fullColor, fullPath, funcName bool
	var layout string

	records := make(chan *Record, uint(2048))
	loggerDefaultTest := newLoggerWithRecords(records)
	defer loggerDefaultTest.Close()

	color = true
	fullPath = true
	c := generateNewConsoleWriterWithOptions(LevelFlagDebug, color, fullColor)
	var name = "console full path"
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
	loggerDefaultTest.Debug("go-log by %s", name)
	loggerDefaultTest.Common("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by %s", name)
	loggerDefaultTest.Transaction("go-log by %s", name)
	loggerDefaultTest.Error("go-log by %s", name)
	loggerDefaultTest.Access("go-log by %s", name)
	loggerDefaultTest.Access("%#v", loggerDefaultTest)
}

func Test_NewConsoleWriterWithFuncName(t *testing.T) {
	var color, fullColor, fullPath, funcName bool
	var layout string

	records := make(chan *Record, uint(2048))
	loggerDefaultTest := newLoggerWithRecords(records)
	defer loggerDefaultTest.Close()

	color = true
	funcName = true
	c := generateNewConsoleWriterWithOptions(LevelFlagDebug, color, fullColor)
	var name = "console func name"
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
	loggerDefaultTest.Debug("go-log by %s", name)
	loggerDefaultTest.Common("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by %s", name)
	loggerDefaultTest.Transaction("go-log by %s", name)
	loggerDefaultTest.Error("go-log by %s", name)
	loggerDefaultTest.Access("go-log by %s", name)
	loggerDefaultTest.Access("%#v", loggerDefaultTest)
}

func Test_NewConsoleWriterWithLayout(t *testing.T) {
	var color, fullColor, fullPath, funcName bool
	var layout string

	records := make(chan *Record, uint(2048))
	loggerDefaultTest := newLoggerWithRecords(records)
	defer loggerDefaultTest.Close()

	color = true
	layout = "20060102T150405.000-0700"
	c := generateNewConsoleWriterWithOptions(LevelFlagDebug, color, fullColor)
	var name = "console layout"
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
	loggerDefaultTest.Debug("go-log by %s", name)
	loggerDefaultTest.Common("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by %s", name)
	loggerDefaultTest.Transaction("go-log by %s", name)
	loggerDefaultTest.Error("go-log by %s", name)
	loggerDefaultTest.Access("go-log by %s", name)
	loggerDefaultTest.Access("%#v", loggerDefaultTest)
}

func Benchmark_NewConsoleWriter(b *testing.B) {
	var color, fullColor, fullPath, funcName bool
	var layout string

	records := make(chan *Record, uint(2048))
	loggerDefaultTest := newLoggerWithRecords(records)
	loggerDefaultTest.SetLevel(DEBUG)
	defer loggerDefaultTest.Close()

	c := generateNewConsoleWriterWithOptions(LevelFlagDebug, color, fullColor)
	var name = "console benchmark test"
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
	loggerDefaultTest.Debug("go-log by %s", name)
	loggerDefaultTest.Common("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by %s", name)
	loggerDefaultTest.Transaction("go-log by %s", name)
	loggerDefaultTest.Error("go-log by %s", name)
	loggerDefaultTest.Access("go-log by %s", name)
	loggerDefaultTest.Access("%#v", loggerDefaultTest)
}

func Benchmark_NewConsoleWriterAll(b *testing.B) {
	var color, fullColor, fullPath, funcName bool
	var layout string

	records := make(chan *Record, uint(2048))
	loggerDefaultTest := newLoggerWithRecords(records)
	loggerDefaultTest.SetLevel(DEBUG)
	defer loggerDefaultTest.Close()
	color = true
	fullColor = true
	fullPath = true
	funcName = true
	layout = "2006-01-02 15:04:05"

	c := generateNewConsoleWriterWithOptions(LevelFlagDebug, color, fullColor)
	var name = "console benchmark test"
	generateRegisterConsoleWriter(loggerDefaultTest, c, fullPath, funcName, layout)
	loggerDefaultTest.Debug("go-log by %s", name)
	loggerDefaultTest.Common("go-log by %s", name)
	loggerDefaultTest.Abnormal("go-log by %s", name)
	loggerDefaultTest.Transaction("go-log by %s", name)
	loggerDefaultTest.Error("go-log by %s", name)
	loggerDefaultTest.Access("go-log by %s", name)
	loggerDefaultTest.Access("%#v", loggerDefaultTest)
}
