/*Package goprobe ...*/
/*
This is used to prepare and send the transaction data to IA
*/
package goprobe

import (
	"context"
	"strconv"

	"github.com/CA-APM/goprobe/internal/utils"
)

//The Transaction type contains the transaction data to be used while
//sending the data to IA
type Transaction struct {
	seqNo    int
	tid      string
	corId    string
	hostName string
	hostPort string
	url      string
	funcName string
}

//EndTransaction is called at the end of every http function handler, this prepares the
//end transaction data in ARF format and send to IA
func (t *Transaction) EndTransaction() {
	defer handlePanic()
	var prmsString string = ",\"prms\":{"
	prmsString += "}"

	m.Lock()
	_, ok := apmTidData[t.tid]
	if ok {
		delete(apmTidData, t.tid)
	}
	m.Unlock()

	arfJSON := utils.ConvertMessageToARF(t.tid, "fnR", t.funcName, t.seqNo, prmsString)
	if getAppInstance().dataConn && ok {
		dataChannel <- arfJSON
	} else {
		getAppInstance().GetLogger().Error("Failed to send EndTransaction data::", map[string]interface{}{
			"data": arfJSON,
		})
	}
}

//StartTransaction is called at the beginning of every http function handler , this prepares the
//transaction data in ARF format and send to IA
//it requires request context and transaction url name
//returns a transaction object
func StartTransaction(ctx context.Context, httpUrl string) Transaction {
	defer handlePanic()
	t := Transaction{
		seqNo:    ctx.Value("seqid").(int) + 1,
		tid:      ctx.Value("tid").(string),
		corId:    ctx.Value("corid").(string),
		hostName: getAppInstance().GetConfig().HostName,
		hostPort: strconv.Itoa(getAppInstance().GetConfig().HostPort),
		url:      httpUrl,
		funcName: "http" + "." + ctx.Value("httpMethod").(string),
	}
	m.Lock()
	apmTidData[t.tid].seqNo = t.seqNo
	m.Unlock()
	var prmsString string = ",\"prms\":{"
	prmsString += "\"url\":\""
	prmsString += t.url
	prmsString += "\",\"hostName\":\""
	prmsString += t.hostName
	prmsString += "\",\"hostPort\":\""
	prmsString += t.hostPort
	prmsString += "\",\"tid\":\""
	prmsString += t.tid
	prmsString += "\",\"corId\":\""
	prmsString += t.corId
	prmsString += "\"}"

	arfJSON := utils.ConvertMessageToARF(t.tid, "fnC", t.funcName, t.seqNo, prmsString)
	if getAppInstance().dataConn {
		dataChannel <- arfJSON
	} else {
		m.Lock()
		_, ok := apmTidData[t.tid]
		if ok {
			delete(apmTidData, t.tid)
		}
		m.Unlock()
		getAppInstance().GetLogger().Error("Failed to send StartTransaction data::", map[string]interface{}{
			"data": arfJSON,
		})
	}

	return t
}
