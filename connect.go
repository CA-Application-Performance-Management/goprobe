/*Package goprobe ...*/
/*
This package is used to establish TCP connection with IA, it uses the configuration
details provide in the config.json file
*/
package goprobe

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CA-APM/goprobe/internal/metric"
	"github.com/CA-APM/goprobe/internal/utils"
)

//TODO now Global need to move
var dataChannel = make(chan string, 1000000)  //buffered dataChannel
var errorChannel = make(chan string, 1000000) //buffered errorChannel
var apmTidData = make(map[string]*seqGen)     //Map to hold the transaction data
var m sync.RWMutex                            //Mutex for synchronization

//The connection type acts like an interface and provided the base
//base methods to be implemented
type connection interface {
	connect()
	read() string
	write(string)
}

//commandconnection type is derived from the base connection type
//This is used to establish the command connection with IA
type commandconnection struct {
	tcpConn *net.TCPConn
}

//This function retrieve the active command connection
func (cmd *commandconnection) getConnection() *net.TCPConn {
	return cmd.tcpConn
}

//This function establish command connection with IA by using the configuration details
//in config.json
func (cmd *commandconnection) connect() {
	service := getAppInstance().GetConfig().HostName + ":" + strconv.Itoa(getAppInstance().GetConfig().HostPort)
	getAppInstance().GetLogger().Debug("commandconnection::Connect", map[string]interface{}{
		"Connection Details": service,
	})
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	validateError(err)
	cmdConn, err := net.DialTCP("tcp", nil, tcpAddr)
	validateError(err)
	//store the connection instance
	cmd.tcpConn = cmdConn
}

//read function is used to read data from the socket
func (cmd *commandconnection) read() string {
	result := make([]byte, 1024)
	bytes, _ := cmd.tcpConn.Read(result)
	getAppInstance().GetLogger().Debug("commandconnection::Read", map[string]interface{}{
		"data": string(result[:bytes]),
	})
	return string(result)
}

//write function is used to write the data to the socket
func (cmd *commandconnection) write(data string) {
	getAppInstance().GetLogger().Debug("commandconnection::Write", map[string]interface{}{
		"data": data,
	})
	_, _ = cmd.tcpConn.Write([]byte(data))
}

//This function retrieves the command connection string in JSON format
func (cmd *commandconnection) getCmdConnection() string {
	probeName := "go"
	probeVersion := 1
	processID := os.Getpid()
	appName := getAppInstance().GetConfig().AppName

	var prmsString string
	prmsString += "{"
	prmsString += "\"appName\":\""
	prmsString += "Go APP\""
	prmsString += ",\"hostName\":\""
	prmsString += getAppInstance().GetConfig().HostName
	prmsString += "\""
	prmsString += "}"
	cmdConnStr := fmt.Sprintf("{\"op\":\"commandConn\",\"probe\":\"%s\",\"ver\":%d,\"instid\":\"%d\",\"pgm\":\"%s\",\"prms\":%s} \r\n", probeName, probeVersion, processID, appName, prmsString)
	getAppInstance().GetLogger().Debug("getCmdConnection", map[string]interface{}{
		"CmdConnection": cmdConnStr,
	})
	return cmdConnStr
}

//dataconnection type is derived from the base connection type
//This is used to establish the Data connection with IA
type dataconnection struct {
	tcpConn *net.TCPConn
}

//This function retrieve the active data connection
func (data *dataconnection) getConnection() *net.TCPConn {
	return data.tcpConn
}

//This function establish data connection with IA by using the configuration details
//in config.json
func (data *dataconnection) connect() {
	service := getAppInstance().GetConfig().HostName + ":" + strconv.Itoa(getAppInstance().GetConfig().HostPort)
	getAppInstance().GetLogger().Debug("dataconnection::Connect", map[string]interface{}{
		"Connection Details": service,
	})
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	validateError(err)
	dataConn, err := net.DialTCP("tcp", nil, tcpAddr)
	validateError(err)
	//store the data connection instance
	data.tcpConn = dataConn
}

//read function is used to read data from the socket
func (data *dataconnection) read() string {
	result := make([]byte, 1024)
	bytes, _ := data.tcpConn.Read(result)
	getAppInstance().GetLogger().Debug("dataconnection::Read", map[string]interface{}{
		"data": string(result[:bytes]),
	})
	return string(result)
}

//write function is used to write the data to the socket
func (data *dataconnection) write(str string) {
	getAppInstance().GetLogger().Debug("dataconnection::Write", map[string]interface{}{
		"data": str,
	})
	_, err := data.tcpConn.Write([]byte(str))
	validateError(err)
}

//This function retrieves the data connection string in JSON format
func (data *dataconnection) getDataConnection() string {
	probeName := "go"
	probeVersion := 1
	processID := os.Getpid()
	dataConnStr := fmt.Sprintf("{\"op\":\"dataConn\",\"probe\":\"%s\",\"ver\":%d,\"instid\":\"%d\",\"tid\":\"%d\",\"pid\":\"%d\"} \r\n", probeName, probeVersion, processID, processID, processID)
	getAppInstance().GetLogger().Debug("getCommandConnection", map[string]interface{}{
		"CmdConnection": dataConnStr,
	})
	return dataConnStr
}

//InitGoProbe estabalishes both connections and runs an error channel
//TODO time lapse needs to be handled in different manner other than sleep
func InitGoProbe() {
	getAppInstance().GetLogger().Debug("InitGoProbe", nil)
	go runErrorChannel()
	go runCommandConnection()
	time.Sleep(5 * time.Second) //Time delay of 5 seconds to make the cmd connection finish before data connection else connection gets terminated
	go runDataConnection()
	go runtimeMetricsTicker()
}

//runDataConnection establishes data connection with IA
//Run's under a go routine
//runs in an infinite loop and looks for any data in dataChannel and if any send it to IA
func runDataConnection() {
	defer handlePanic()
	getAppInstance().GetLogger().Debug("runDataConnection", nil)
	//Establish dataconnection
	dataConn := dataconnection{}
	dataConn.connect()
	defer dataConn.getConnection().Close()
	//Write data connection string
	dataConnectionStr := dataConn.getDataConnection()
	dataConn.write(dataConnectionStr)
	//infinite loop that reads data from dataChannel
	for {
		msg := <-dataChannel

		getAppInstance().GetLogger().Debug("runDataConnection", map[string]interface{}{
			"Data to Agent": msg,
		})
		//write to data connection socket
		dataConn.write(msg)
	}
}

//runCommandConnection establishes command connection with IA
//run's under a go routine
//run's in an infinite loop and looks for "speak" call from IA and send a response to make the connection "alive"
func runCommandConnection() {
	defer handlePanic()
	getAppInstance().GetLogger().Debug("runCommandConnection", nil)
	//Establish command connection
	cmdConn := commandconnection{}
	cmdConn.connect()
	defer cmdConn.getConnection().Close()
	//write command connection string
	commandConnectionStr := cmdConn.getCmdConnection()
	cmdConn.write(commandConnectionStr)
	//infinite loop that makes the command connection Alive
	for {
		result := cmdConn.read()
		if len(result) == 0 {
			continue
		}
		//check if "speak" is receive from IA, then send the response to make the connection alive
		if strings.Contains(string(result), "speak") {
			ccSpeakStr := fmt.Sprintf("{\"op\":\"arf\"}\n")
			//write to command connection socket
			cmdConn.write(ccSpeakStr)
			continue
		}
		//TODO to verify the config message
		if strings.Contains(string(result), "config") {
			continue
		}
	}
}

//validateError checks for any incoming error and send it to errorChannel for processing
func validateError(err error) {
	if err != nil {
		errorChannel <- err.Error()
		getAppInstance().dataConn = false
		getAppInstance().cmdConn = false
	}
}

//This function runs under go routine with an infinite loop looking for any errors in error channel
//if we encounter any error in the channel respective action will be taken
func runErrorChannel() {
	defer handlePanic()
	getAppInstance().GetLogger().Debug("runErrorChannel", nil)
	for {
		val := <-errorChannel
		getAppInstance().GetLogger().Error("Runtime Error", map[string]interface{}{
			"Error": val,
		})
	}
}

//This function will spin a new goRoutine based on the time interval mentioned in the config.json
// Before running a ticker it will also collect the static metrics and send to IA using dataChannel
// A new goroutine will collect the runtime metrics and push the metrics to IA using dataChannel
func runtimeMetricsTicker() {
	defer handlePanic()
	getAppInstance().GetLogger().Debug("runtimeMetricsTicker", map[string]interface{}{
		"Runtime Interval": strconv.Itoa(getAppInstance().GetConfig().Interval),
	})
	publishStaticMetrics(metric.CollectStaticMetrics())
	ticker := time.NewTicker(time.Duration(getAppInstance().GetConfig().Interval) * time.Second)
	go func() {
		for range ticker.C {
			defer handlePanic()
			metrics := metric.CollectRuntimeMetrics()
			publishRuntimeMetrics(metrics)
		}
	}()
}

//This function will iterate through the collected metrics and push to the dataChannel
func publishRuntimeMetrics(metrics []metric.Metric) {
	for _, metric := range metrics {
		getAppInstance().GetLogger().Debug("publishRuntimeMetrics", map[string]interface{}{
			"Metric": utils.ConvertMetricToARF(metric.Op, metric.Id, metric.Type, metric.Value),
		})
		if getAppInstance().dataConn {
			dataChannel <- utils.ConvertMetricToARF(metric.Op, metric.Id, metric.Type, metric.Value)
		} else {
			getAppInstance().GetLogger().Error("Failed to send RuntimeMetrics data::", map[string]interface{}{
				"data": utils.ConvertMetricToARF(metric.Op, metric.Id, metric.Type, metric.Value),
			})
		}
	}
}

//This function will iterate through the collected static metrics and push to the dataChannel
func publishStaticMetrics(metrics []metric.Metric) {
	for _, metric := range metrics {
		getAppInstance().GetLogger().Debug("publishStaticMetrics", map[string]interface{}{
			"Metric": utils.ConvertMetricToARF(metric.Op, metric.Id, metric.Type, metric.Value),
		})
		if getAppInstance().dataConn {
			dataChannel <- utils.ConvertMetricToARF(metric.Op, metric.Id, metric.Type, metric.Value)
		} else {
			getAppInstance().GetLogger().Error("Failed to send StaticMetrics data::", map[string]interface{}{
				"data": utils.ConvertMetricToARF(metric.Op, metric.Id, metric.Type, metric.Value),
			})
		}
	}
}
