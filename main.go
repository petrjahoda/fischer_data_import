package main

import (
	"github.com/kardianos/service"
	"time"
)

const version = "2020.4.1.29"
const serviceName = "Fischer Data Import Service"
const serviceDescription = "Download users, products and product groups from Fischer mssql database and imports them into Zapsi mysql database"
const downloadInSeconds = 600

var serviceRunning = false
var processRunning = false
var zapsiConfig = "zapsi_uzivatel:zapsi@tcp(zapsidatabase:3306)/zapsi2?charset=utf8&parseTime=True&loc=Local"
var fischerConfig = "sqlserver://zapsi_read:fis32zap@10.193.21.37:1444/FPCZ?database=HELIOS_fpcz"

type program struct{}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	logInfo("MAIN", serviceDescription)
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	program := &program{}
	s, err := service.New(program, serviceConfig)
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
}

func (p *program) Start(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	serviceRunning = true
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	serviceRunning = false
	if processRunning {
		logInfo("MAIN", serviceName+" ["+version+"] stopping...")
		time.Sleep(1 * time.Second)
	}
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

func (p *program) run() {
	for serviceRunning {
		processRunning = true
		start := time.Now()
		logInfo("MAIN", serviceName+" ["+version+"] running")
		importData()
		sleepTime := downloadInSeconds*time.Second - time.Since(start)
		logInfo("MAIN", "Sleeping for "+sleepTime.String())
		time.Sleep(sleepTime)
		processRunning = false
	}
}
