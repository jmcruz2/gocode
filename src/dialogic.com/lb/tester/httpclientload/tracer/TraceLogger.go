package tracer

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LoggerDebugMsg struct {
	ClientId                  int
	MsgId                     int
	ReusingConnection         bool
	TotalMsgPerClient         int64
	RoundTripResponseDuration time.Duration
	ContentSizeBytes          int64
}

type LOG_TYPE int

const LOG_TYPE_10SEC_MSG LOG_TYPE = 1
const LOG_TYPE_DEBUG LOG_TYPE = 2
const LOG_TYPE_MAX_MSG LOG_TYPE = 3

type LoggerMsg struct {
	logType             LOG_TYPE
	numClients          int
	totalMsgPer10Sec    int64
	Marker              string
	Devider             bool
	LineFeed            bool
	Attributes          map[string]time.Duration
	pingPong            string
	avgContentSizeBytes int64
	maxNumberOfMsgs     int64

	debugStruct LoggerDebugMsg
}

type TraceLogger struct {
	DontExit   bool
	lbChannel  chan LoggerMsg
	fLog       *os.File
}

func (tl *TraceLogger) Init(bcSize int) {
	//fmt.Println("TLOGGER Init to buffer size: ", bcSize)
	tl.lbChannel = make(chan LoggerMsg, bcSize)
	tl.DontExit = true
	var err error

	fileName := "httpLbTestResult-"

	const fileTiimeLayout = "2006-01-02-15:04:05"
	fileName += time.Now().Format(fileTiimeLayout)

	tl.fLog, err = os.OpenFile("lb-test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		fmt.Println("error opening file: %v", err)
	}
	log.SetOutput(tl.fLog)
	//log.SetPrefix(time.Now().Format("15:04:05") )
	log.SetFlags(log.Ltime)
	//log.Println("using log Println TraceLogger")

	go tl.run(bcSize)
	//fmt.Println("TLOGGER return from go run")

}

func (tl *TraceLogger) run(bcSize int) {
	//fmt.Println("TLOGGER Enter run")

	for tl.DontExit {
		//fmt.Println("TLOGGER: waiting for logTrace Msg from lbChannel")

		traceMsg := <-tl.lbChannel //get from channel queue
		//fmt.Println("TLOGGER: got logTrace Msg from lbChannel")

		if traceMsg.logType == LOG_TYPE_DEBUG {
			header1 := fmt.Sprintf("%20s%20s%20s%20s%30s%25s", "ClientId", "MsgId", "ReusingConnection", "TotalMsgPerClient", "RoundTripResponseDuration", "ContentSizeBytes")
			log.Println(header1)
			line1 := fmt.Sprintf("%20v%20v%20v%20v%20v%25v\n", traceMsg.debugStruct.ClientId, traceMsg.debugStruct.MsgId, traceMsg.debugStruct.ReusingConnection, traceMsg.debugStruct.TotalMsgPerClient,
				traceMsg.debugStruct.RoundTripResponseDuration, traceMsg.debugStruct.ContentSizeBytes)
			log.Println(line1)

		} else if traceMsg.logType == LOG_TYPE_10SEC_MSG {
			if traceMsg.Devider {
				log.Println("===============================================================================================================================")
			}

			totalMsgPerSec := traceMsg.totalMsgPer10Sec / 10
			header1 := fmt.Sprintf("%20s%20s%20s%20s%20s%25s", "Buffer Type", "Clients", "ReqPer/10Secs", "ReqPer/sec", "Avg Response Time", "Peak Response Time ")
			log.Println(header1)
			line1 := fmt.Sprintf("%20v%20v%20v%20v%20v%25v\n", traceMsg.pingPong, traceMsg.numClients, traceMsg.totalMsgPer10Sec, totalMsgPerSec,
				traceMsg.Attributes["AvgRoundTripResponseDuration"], traceMsg.Attributes["PeakRoundTripResponseDuration"])
			log.Println(line1)

			if traceMsg.LineFeed {
				log.Println()
			}
		} else if traceMsg.logType == LOG_TYPE_MAX_MSG {
			log.Println("================================Total Test Run Values  =================================================================================")
			log.Println("Total Messages: ", traceMsg.maxNumberOfMsgs)
			fmt.Println("Total Messages: ", traceMsg.maxNumberOfMsgs)
			log.Println("End of Load Test")
		}

	}

	tl.fLog.Close()

}

func (tl *TraceLogger) Log(traceMsg LoggerMsg) {

	tl.lbChannel <- traceMsg //write to channel queue
}
