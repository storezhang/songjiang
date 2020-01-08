package common

import (
    "songjiang/sign"
)

// Config 程序整体配置
type Config struct {
    Songjiang Songjiang
    Hao4k     sign.Hao4k
    Apps      []App
}

// Songjiang 程序整体配置
type Songjiang struct {
    Debug          bool   `default:"false"`
    LogLevel       string `default:"info" yaml:"logLevel" toml:"logLevel"`
    Chans          []ServerChan
    Template       Template
    BrowserWidth   int    `default:"1920" yaml:"browserWidth" toml:"browserWidth"`
    BrowserHeight  int    `default:"1080" yaml:"browserHeight" toml:"browserHeight"`
    BrowserTimeout string `default:"30m" yaml:"browserTimeout" toml:"browserTimeout"`
    Redo           string `default:"5m"`
    RetryLimit     uint   `default:"30" yaml:"retryLimit" toml:"retryLimit"`
}

// App 应用配置
type App struct {
    Name      string `default:"应用1"`
    Chans     []ServerChan
    Template  Template
    Type      string `default:"hao4k"`
    Cookies   string
    StartTime string `default:"8:00" yaml:"startTime" toml:"startTime"`
    EndTime   string `default:"23:00" yaml:"endTime" toml:"endTime"`
}

// Template 模板配置
// 用于：推送
type Template struct {
    Title   string `default:"'签到后：{{.Result.After}}，签到前{{.Result.Before}}'"`
    Content string `default:"'任务名称：{{.App.Name}} {{.Result.Msg}}'"`
}
