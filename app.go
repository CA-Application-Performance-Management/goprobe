/*Package goprobe ...*/
/*
This package the main package exposed to the external world and all the interfaces
accessed will be added in this package
*/
package goprobe

import (
	"sync"

	cc "github.com/CA-APM/goprobe/internal/config"
	ll "github.com/CA-APM/goprobe/internal/logger"
)

//application type hold the app instance of the whole package
type application struct {
	config   cc.Configuration
	logger   ll.Logger
	dataConn bool
	cmdConn  bool
}

//GetConfig retrieve the application configuration instance
func (app *application) GetConfig() cc.Configuration {
	return app.config
}

//GetLogger retrieves the application logger instance
func (app *application) GetLogger() ll.Logger {
	return app.logger
}

var instance *application
var once sync.Once

//getAppInstance creates a singleton instance of application type
//this instance is used by other packages
func getAppInstance() *application {
	once.Do(func() {
		myConfig, myLog := cc.ReadConfig()
		instance = &application{
			config:   myConfig,
			logger:   myLog,
			dataConn: true,
			cmdConn:  true,
		}
	})
	return instance
}
