package main

import (
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/tcpless"
	"github.com/dimonrus/tcpless/example/json_client/config"
	"os"
	"os/signal"
)

var App gocli.Application

var cliArgs = gocli.Arguments{"app": gocli.Argument{
	Type:  "string",
	Label: "Application type",
	Name:  "app",
}}

func InitApp() {
	var cfg config.Config
	configPath, _ := gocli.DNApp{}.GetAbsolutePath("config", "example/json_client")
	App = gocli.NewApplication(os.Getenv("ENV"), configPath, &cfg)
	App.SetLogger(gocli.NewLogger(gocli.LoggerConfig{}))
	App.ParseFlags(&cliArgs)
}

func main() {
	InitApp()
	var server *tcpless.Server
	switch cliArgs["app"].GetString() {
	case "server":
		server = StartServer(&App.GetConfig().(*config.Config).TCPLess, App)
		sig := make(chan os.Signal, 1)
		go func() {
			signal.Notify(sig, os.Interrupt)
		}()
		<-sig
		server.Stop()
	case "client":
		StartClient(&App.GetConfig().(*config.Config).TCPLess, App)
	}
	App.GetLogger().Warnln("shutdown")
}
