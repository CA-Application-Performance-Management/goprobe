/*Package goprobe ...*/
/*
This is used to prepare and send the segment data to IA
*/
package goprobe

import (
	"context"
	"math/rand"
	"strconv"
	"strings"

	"github.com/CA-APM/goprobe/internal/utils"
)

const (
	fileType   = "fs"
	sqlType    = "mysql"
	oracleType = "oracle"
)

//The Segment type contains the segment data to be used while
//sending the data to IA
type Segment struct {
	seqNo    int
	tid      string
	funcName string
	segType  string
}

//EndSegment is called at the end of every function, this prepares the
//end segment data in ARF format and send to IA
func (s *Segment) EndSegment() {
	defer handlePanic()
	var prmsString string = ",\"prms\":{"
	prmsString += "}"

	m.Lock()
	_, ok := apmTidData[s.tid]
	if ok {
		apmTidData[s.tid].seqNo = apmTidData[s.tid].seqNo - 1
	}
	m.Unlock()

	arfJSON := utils.ConvertMessageToARF(s.tid, "fnR", s.segType+"."+s.funcName, s.seqNo, prmsString)
	if getAppInstance().dataConn && ok {
		dataChannel <- arfJSON
	} else {
		getAppInstance().GetLogger().Error("Failed to send EndSegment data::", map[string]interface{}{
			"data": arfJSON,
		})
	}
}

//StartSegment is called at the beginning of every function, this prepares the
//segment data in ARF format and send to IA
//it requires request context and segment name
//returns a segment object
func StartSegment(ctx context.Context, name string, optional ...string) Segment {
	defer handlePanic()
	var prmsString string = ""
	var s Segment
	var sType string
	ok := true

	if len(optional) > 0 && len(optional) == 2 {
		var info dbInfo
		//retrieve the db details and form the prms based on the details
		info.getDbdetails(optional[0])
		sType = info.dbType
		prmsString = ",\"prms\":{"
		prmsString += "\"url\":\""
		prmsString += optional[0]
		if len(info.dbName) > 0 {
			prmsString += "\",\"dbname\":\""
			prmsString += info.dbName
		}
		if len(info.dbHost) > 0 {
			prmsString += "\",\"host\":\""
			prmsString += info.dbHost
		}
		if len(info.dbPort) > 0 {
			prmsString += "\",\"port\":\""
			prmsString += info.dbPort
		}
		if len(info.dbType) > 0 {
			prmsString += "\",\"database\":\""
			prmsString += info.dbType
		}
		prmsString += "\",\"query\":\""
		prmsString += optional[1]
		prmsString += "\",\"sql\":\""
		prmsString += optional[1]
		prmsString += "\",\"traceId\":\""
		prmsString += ctx.Value("tid").(string)
		prmsString += "\",\"corId\":\""
		prmsString += ctx.Value("corid").(string)
		prmsString += "\"}"
	} else {
		sType = fileType
	}
	if ctx != nil {
		tid := ctx.Value("tid").(string)
		seqNo := 1
		m.Lock()
		_, ok = apmTidData[ctx.Value("tid").(string)]
		if ok {
			seqNo = apmTidData[ctx.Value("tid").(string)].seqNo + 1
			apmTidData[tid].seqNo = seqNo
		}
		m.Unlock()
		s = Segment{
			seqNo:    seqNo,
			tid:      tid,
			funcName: name,
			segType:  sType,
		}
	} else {
		s = Segment{
			seqNo:    1,
			tid:      strconv.FormatInt(rand.Int63(), 10),
			funcName: name,
			segType:  sType,
		}
	}

	arfJSON := utils.ConvertMessageToARF(s.tid, "fnC", s.segType+"."+s.funcName, s.seqNo, prmsString)
	if getAppInstance().dataConn && ok {
		dataChannel <- arfJSON
	} else {
		getAppInstance().GetLogger().Error("Failed to send StartSegment data::", map[string]interface{}{
			"data": arfJSON,
		})
	}
	return s
}

//The dbInfo type contains the information about the database details
//this is used to send the database backeneds to EM
type dbInfo struct {
	dbType string
	dbName string
	dbHost string
	dbPort string
}

//getDbdetails will retrieve the database information using the db connection stringlike
// dbType -- Type of the database like mysql, oracle
// dbName -- Name of the database
// dbHost -- Host where the database exists
// dbPort -- Port used to connect to the database
func (db *dbInfo) getDbdetails(connString string) {
	if strings.Contains(connString, "Data Source") {
		db.dbType = oracleType
	} else {
		db.dbType = sqlType
		db.dbPort = "3306" //default port for mysql
		dbSlice := strings.Split(connString, ";")
		if dbSlice[0] != connString {
			for _, detail := range dbSlice {
				result := strings.Split(detail, "=")
				if result[0] != detail {
					key, value := result[0], result[1]
					if key == "Database" {
						db.dbName = value
					} else if key == "Server" {
						db.dbHost = value
					} else if key == "Port" {
						db.dbPort = value
					}
				}
			}
		}
	}
}
