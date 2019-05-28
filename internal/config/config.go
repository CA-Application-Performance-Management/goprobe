/*Package config ...*/
/*
This package reads the user configurations and if any details missed from the configurations
default values are used
*/
package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/CA-APM/goprobe/internal/logger"
	util "github.com/CA-APM/goprobe/internal/utils"
)

/*Configuration ...*/
/*Contains all the details that the SDK uses and can be configured by the user in the file config.json
AppName -- This is the application name that will be displayed in the EM/ATC, if this is not configured default name 'goProbe' will be taken
HostName -- This is the IA HostName and this shoule be a valid one, if invalid error is logged and process exits
HostPort -- This is the IA HostPort and this shoule ba valid one, if invalid error is logged and process exits
LogPath -- A valid path that includes the filename(error.log), if invalid the log is redirected to standard output
DebugEnabled -- true/false, this enable to write the debugging logs
Interval -- This is the runtime interval to capture the runtime metrics by the SDK, if not configured default 120 seconds is used
*/
type Configuration struct {
	AppName      string
	HostName     string
	HostPort     int
	LogPath      string
	DebugEnabled bool
	Interval     int
}

//GetApplicationName function gets the AppName of the configuration
func (c *Configuration) GetApplicationName() string {
	return c.AppName
}

//SetApplicationName function sets the AppName of the configuration
func (c *Configuration) SetApplicationName(appName string) {
	c.AppName = appName
}

//GetHostName function gets the HostName of the configuration
func (c *Configuration) GetHostName() string {
	return c.HostName
}

//SetHostName function sets the HostName of the configuration
func (c *Configuration) SetHostName(hostName string) {
	c.HostName = hostName
}

//GetHostPort function gets the HostPort of the configuration
func (c *Configuration) GetHostPort() int {
	return c.HostPort
}

//SetHostPort function sets the HostPort of the configuration
func (c *Configuration) SetHostPort(hostPort int) {
	c.HostPort = hostPort
}

//GetLogPath function gets the LogPath of the configuration
func (c *Configuration) GetLogPath() string {
	return c.LogPath
}

//SetLogPath function sets the LogPath of the configuration
func (c *Configuration) SetLogPath(logPath string) {
	c.LogPath = logPath
}

//GetDebugEnabled function gets the DebugEnabled of the configuration
func (c *Configuration) GetDebugEnabled() bool {
	return c.DebugEnabled
}

//SetDebugEnabled function sets the DebugEnabled of the configuration
func (c *Configuration) SetDebugEnabled(debugEnabled bool) {
	c.DebugEnabled = debugEnabled
}

//GetInterval function gets the Interval of the configuration
func (c *Configuration) GetInterval() int {
	return c.Interval
}

//SetInterval function sets the Interval of the configuration
func (c *Configuration) SetInterval(interval int) {
	c.Interval = interval
}

const configFile string = "config.json"

var goLog logger.Logger

//ReadConfig is used to read the configurations and sets the logger for the SDK
func ReadConfig() (Configuration, logger.Logger) {
	conf := loadConfig()
	loadLogger(&conf)
	validate(&conf)
	return conf, goLog
}

//loadConfig reads the config.json and decode the configurations
//Reads the config file and exits the process if no file/failed to decode
func loadConfig() Configuration {
	confPath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "CA-APM", "goprobe", "config.json")
	file, err := os.Open(confPath)
	if err != nil {
		log.Fatal("unable to open config file: ", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatal("unable to decode config JSON: ", err)
	}
	return conf
}

//loadLogger reads the log info from config file and sets the logger
//If valid log file is present, logger sets to the log file
//If invalid, logger sets to the standard output
func loadLogger(c *Configuration) {

	if len(c.GetLogPath()) > 0 && util.IsValidPath(filepath.Dir(c.GetLogPath())) {
		file, err := os.OpenFile(c.GetLogPath(), os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		goLog = logger.New(file, c.GetDebugEnabled())
	} else {
		goLog = logger.New(os.Stdout, c.GetDebugEnabled())
	}
}

//validate is used to check valid configurations are present, if not present set the defaults
//checks for application name if not present default "goProbe" is used
//checks for internal if not present default "120" seconds is used
//checks for HostName if not valid exis the process
func validate(c *Configuration) {

	if len(c.GetApplicationName()) == 0 {
		c.SetApplicationName("goProbe")
		goLog.Info("Missing Application Name", map[string]interface{}{
			"appName": "",
		})
	}
	if c.GetInterval() == 0 {
		c.SetInterval(120)
		goLog.Info("Missing Runtime Interval", map[string]interface{}{
			"interval": "",
		})
	}
	service := c.GetHostName() + ":" + strconv.Itoa(c.GetHostPort())
	if util.IsValidHost(service) {
		goLog.Debug("Host Details", map[string]interface{}{
			"HostName": c.GetHostName(),
			"HostPort": c.GetHostPort(),
		})
	} else {
		goLog.Error("Invalid Host Details", map[string]interface{}{
			"HostName": c.GetHostName(),
			"HostPort": c.GetHostPort(),
		})
		os.Exit(1)
	}

}
