# go-log

## Introduction
a useful go log, support multiple writer.

## TODO

- [x] base logger
- [x] console writer
- [x] file writer
- [x] log file rotate
- [ ] kafka writer (soon)
- [ ] sensitive information protection (soon)

## ENV

The go version shall >= `1.16`

## Install

```text
go get -u -v github.com/legofun/go-log
```

## Usage

```go
import github.com/legofun/go-log

...

golog.SetupLog(golog.LogConfig{
    Level:    "debug",
    Debug:    false,
    FullPath: true,
    ConsoleWriter: golog.ConsoleWriterOptions{
        Enable: true,
        Color:  true,
        //Level:  "abnormal",
    },
    FileWriter: golog.FileWriterOptions{
        //Level:      "",
        Filename: "./test/golog-test-%Y%M%D%H%m.log",
        Enable:   true,
        Rotate:   true,
        Daily:    false,
        Hourly:   false,
        Minutely: true,
        //MaxDays:    0,
        //MaxHours:   0,
        //MaxMinutes: 0,
    },
})

golog.Debug("this is debug log")
golog.Common("this is common log")
golog.Abnormal("this is abnormal log")
golog.Transaction("this is transaction log")
golog.Error("this is error log")
golog.Access("this is access log")
```

## License

Use of go-log is governed by the MIT License
