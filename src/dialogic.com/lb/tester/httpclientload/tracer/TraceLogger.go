package tracer

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LoggerMsg struct {
	numClients	 int
	totalMsgPer10Sec uint64
	Marker           string
	Devider          bool
	LineFeed         bool
	Attributes map[string]time.Duration
	pingPong         string
}

type TraceLogger struct {
	DontExit  bool
	lbChannel chan LoggerMsg
	fLog      *os.File
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

		if traceMsg.Devider {
			log.Println("===============================================================================================================================")
		}

		header1 := fmt.Sprintf("%16s%16s%16s%16s",  "Buffer Type", "Clients", "MSGS/10Secs", "tr/sec" )
		log.Println(header1)
		line1 := fmt.Sprintf("%16v%16v%16v%16v\n", traceMsg.pingPong, traceMsg.numClients, traceMsg.totalMsgPer10Sec, 0 )
		log.Println(line1)

		header2 := fmt.Sprintf("%16s%16s%16s%16s%16s%16s",   "MaxResponse", "MaxServer ", "MaxLatency", "MinResponse", "MinServer", "MinLatency" )
		log.Println(header2)
		line2 := fmt.Sprintf("%16v%16v%16v%16v%16v%16v\n",  traceMsg.Attributes["MaxResponse"], traceMsg.Attributes["MaxServer"],
			traceMsg.Attributes["MaxLatency"], traceMsg.Attributes["MinResponse"], traceMsg.Attributes["MinServer"], traceMsg.Attributes["MinLatency"])
		log.Println(line2)


		header3 := fmt.Sprintf("%16s%16s%16s", "AvgResponse", "AvgServer", "AvgLatency" )
		log.Println(header3)
		line3 := fmt.Sprintf("%16v%16v%16v", traceMsg.Attributes["AvgLatency"], traceMsg.Attributes["AvgResponse"],traceMsg.Attributes["AvgServer"] )
		log.Println(line3)


		//msg := fmt.Sprintf("%16v%16v%16v%16v\n%16v%16v%16v%16v%16v%16v\n%16v%16v%16v", traceMsg.pingPong, traceMsg.numClients, traceMsg.totalMsgPer10Sec, 0, traceMsg.Attributes["MaxResponse"], traceMsg.Attributes["MaxServer"],
		//	traceMsg.Attributes["MaxLatency"], traceMsg.Attributes["MinResponse"], traceMsg.Attributes["MinServer"], traceMsg.Attributes["MinLatency"],
		//	traceMsg.Attributes["AvgLatency"], traceMsg.Attributes["AvgResponse"],traceMsg.Attributes["AvgServer"] )
		//	log.Println(msg)
		//for attrName, attrVal := range traceMsg.Attributes {

		//	log.Println(attrName, attrVal)
		//}

		if traceMsg.LineFeed {
			log.Println()
		}

	}
	//fmt.Println("TLOGGER Exit run")

	tl.fLog.Close()

}

func (tl *TraceLogger) Log(traceMsg LoggerMsg) {
	//fmt.Println("TLOGGER Log : putting tl.lbChannel <- traceMsg" , traceMsg )

	tl.lbChannel <- traceMsg //write to channel queue
}


func (tl *TraceLogger) LogTestRunTotals(totalMsgs uint64) {

	log.Println("================================Total Test Run Values  =================================================================================")
	log.Println("Total Messages: ", totalMsgs)
	fmt.Println("Total Messages: ", totalMsgs)
	log.Println("End of Load Test")

}

func ( t1 *TraceLogger )Debug(dbgMsg string) {
	log.Println("DEBUG: " , dbgMsg)
}

func ( t1 *TraceLogger )DebugUInt64(dbgMsg uint64) {
	log.Println("DEBUG: " , dbgMsg)
}


