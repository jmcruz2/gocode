package tracer

import (
	"fmt"
	"sync"
	"time"
)

type CHANNEL_PINGPONG_TYPE int

type CommandLineParams struct {
	NumOfClients   *int
	CallsPerSecond *int
	RepetionInMin  *int
	LbUri          *string
	DbgFlag        *bool
}

const (
	CHANNEL_TYPE_PING CHANNEL_PINGPONG_TYPE = 100
	CHANNEL_TYPE_PONG CHANNEL_PINGPONG_TYPE = 200
)

type TraceDataSample struct {
	clientID                  int
	msgID                     int
	newConnection             bool
	roundTripresponseDuration time.Duration
	contentSizeBytes          int64
}

//var channelActiveType2 CHANNEL_PINGPONG_TYPE
//var channelPreviousActiveType2 CHANNEL_PINGPONG_TYPE
type TraceEngine struct {
	logger               TraceLogger
	tickerMgr            TraceTickerMgr
	waitOnTraceEntryLock *sync.Mutex
	activeChannel        *chan TraceDataSample
	previousChan         *chan TraceDataSample

	pingChannel chan TraceDataSample
	pongChannel chan TraceDataSample

	dontExit        bool
	dontExitChannel chan bool

	channelActiveType         CHANNEL_PINGPONG_TYPE
	channelPreviousActiveType CHANNEL_PINGPONG_TYPE

	mping pingMatrics
	mpong pongMatrics

	startedFlag bool

	testDone chan bool

	numClients int

	totalMessageProcessDuringTest int64

	work10Lock *sync.Mutex

	DoExit bool

	debugFlag bool
}

func (t *TraceEngine) Ctor(channelBufSize int, reps int, numOfClients int, done chan bool, dbg bool) {
	//fmt.Println("TraceEngine - CTOR")
	t.waitOnTraceEntryLock = &sync.Mutex{}
	t.work10Lock = &sync.Mutex{}
	t.testDone = done
	t.debugFlag = dbg
	t.numClients = numOfClients
	t.pingChannel = make(chan TraceDataSample, channelBufSize)
	t.pongChannel = make(chan TraceDataSample, channelBufSize)
	t.activeChannel = &t.pingChannel
	t.previousChan = &t.pingChannel
	t.dontExitChannel = make(chan bool)
	t.dontExit = true
	t.channelActiveType = CHANNEL_TYPE_PING
	t.channelPreviousActiveType = CHANNEL_TYPE_PING

	t.totalMessageProcessDuringTest = 0
	t.DoExit = false

	t.mping = pingMatrics{Total10SecTraceMatrics{}}
	t.mpong = pongMatrics{Total10SecTraceMatrics{}}
	t.mping.Reset()
	t.mpong.Reset()

	t.tickerMgr = TraceTickerMgr{}
	t.tickerMgr.Ctor(reps)
	t.logger = TraceLogger{}
	t.logger.Init(1000) //1000 message buffer plenty
	t.startedFlag = false

	//tNow := time.Now()
	//t.logger.Debug(fmt.Sprintf("**********************************************************************************"))
	//t.logger.Debug(fmt.Sprintf("New Load Test Started at %v", tNow))
	//t.logger.Debug("Column Definitions:")
	//t.logger.Debug("Clients=>Number of clients making calls - MSGS/10secs=>Number of Messages Sent to server per 10secs")
	//t.logger.Debug("tr/sec=>Max rate calls/sec")
	//
	//t.logger.Debug("MaxResponse: Maximum Total Response in the 10sec period maxResponse=maxLatency+MaxServer")
	//t.logger.Debug("MaxServer:   Maximum time taken by the server to process request within the 10sec period")
	//t.logger.Debug("MaxLatency:  Maximum xmit and receive latency - maxLatency = MaxResponse - MaxServer (time on network)")
	//
	//t.logger.Debug("MinResponse: Min Total Response in the 10sec period maxResponse=maxLatency+MaxServer")
	//t.logger.Debug("MinServer:   Min time taken by the server to process request within the 10sec period")
	//t.logger.Debug("MinLatency:  Min xmit and receive latency - minLatency = MinResponse - MinServer (time on network)")
	//
	//t.logger.Debug("AvgResponse: Average Total Response in the 10sec period AvgResponse=AvgLatency+AvgServer")
	//t.logger.Debug("AvgServer:   Average time taken by the server to process request within the 10sec period")
	//t.logger.Debug("AvgLatency:  Average xmit and receive latency - AvgLatency = AvgResponse - AvgServer (time on network)")

}

func (t *TraceEngine) Start() {
	//fmt.Println("Starting Trace Engine")
	t.run()
}

func (t *TraceEngine) setChannelActive(val CHANNEL_PINGPONG_TYPE) {
	t.channelActiveType = val
}

func (t *TraceEngine) getChannelActive() CHANNEL_PINGPONG_TYPE {
	return t.channelActiveType
}

func (t *TraceEngine) TraceEntry(traceData TraceDataSample) {

	//fmt.Println("TraceEntry() calling waitOnTraceEntryLock LOCK")
	//t.waitOnTraceEntryLock.Lock()

	//fmt.Println("Writing: ", channelActiveType2)
	if t.debugFlag {
		fmt.Println("TraceEntry() calling waitOnTraceEntryLock LOCK")
		if t.channelActiveType == CHANNEL_TYPE_PING {
			fmt.Println("TraceEntry() Writing to pingChannel")
			//t.pingChannel <- traceData
		} else {
			fmt.Println("TraceEntry() Writing to pongChannel")
			//	t.pongChannel <- traceData
		}
	}

	*t.activeChannel <- traceData

	if t.debugFlag {
		fmt.Println("TraceEntry() calling waitOnTraceEntryLock UNLOCK")
	}

	//t.waitOnTraceEntryLock.Unlock()

}

func (t *TraceEngine) calc10SecsLogEntry(tdSample TraceDataSample, chType CHANNEL_PINGPONG_TYPE) {

	if t.debugFlag {
		fmt.Println("calc10SecsLogEntry")
	}
	var tm *Total10SecTraceMatrics

	if chType == CHANNEL_TYPE_PING {
		if t.debugFlag {
			fmt.Println("Writing calc10SecsLogEntry - Receive data from Ping Channel")
		}
		tm = t.mping.GetMatrics()
	} else {
		if t.debugFlag {
			fmt.Println("Writing calc10SecsLogEntry - Receive data from Pong Channel")
		}
		tm = t.mpong.GetMatrics()
	}

	tm.totalMsgsPer10Secs = tm.totalMsgsPer10Secs + 1
	t.totalMessageProcessDuringTest = t.totalMessageProcessDuringTest + 1

	//if tm.avgLatencyCnt == 0 {
	//	if t.debugFlag {
	//		fmt.Println("ZERO DATA")
	//	}
	//	tm.avgLatencyCnt++
	//	tm.avgLatencyTotal = tdSample.latencyDuration
	//	tm.avgLatencyDuration = tdSample.latencyDuration
	//	tm.maxLatencyDuration = tdSample.latencyDuration
	//	tm.minLatencyDuration = tdSample.latencyDuration
	//
	//	tm.avgServerCnt++
	//
	//
	//	tm.avgTotalRespCnt++
	//
	//
	//} else {
	{
		if t.debugFlag {
			fmt.Println("NONE ZERO DATA")
		}

		tm.AvgContentSizeBytes = (tm.AvgContentSizeBytes + tdSample.contentSizeBytes) / tm.totalMsgsPer10Secs

		tm.sumRoundTripResponseDuration = tm.sumRoundTripResponseDuration + tdSample.roundTripresponseDuration

		var d64 int64
		d64 = tm.sumRoundTripResponseDuration.Nanoseconds()
		div64 := d64 / tm.totalMsgsPer10Secs

		//tm.AvgRoundTripResponseDuration = tm.sumRoundTripResponseDuration / tm.totalMsgsPer10Secs
		tm.AvgRoundTripResponseDuration = time.Duration(div64)

		if tm.PeakRoundTripResponseDuration < tdSample.roundTripresponseDuration {
			tm.PeakRoundTripResponseDuration = tdSample.roundTripresponseDuration
		}
	}

}

func (t *TraceEngine) work10SecTimeout() {

	if t.debugFlag {
		fmt.Println("work10SecTimeout - calling work10Lock LOCK")
	}

	t.work10Lock.Lock()
	var traceDataSample TraceDataSample
	var exitLoop = false

	var tm *Total10SecTraceMatrics
	var prevPingPong string

	//stime := time.Now()

	//var bufChan chan TraceDataSample
	if t.channelPreviousActiveType == CHANNEL_TYPE_PING {
		//if channelPreviousActiveType2 == CHANNEL_TYPE_PING {
		prevPingPong = "PING"
		//	bufChan = t.pingChannel
		tm = t.mping.GetMatrics()
		if t.debugFlag {

			fmt.Println("work10secTimeout - reading from PING BUFFER")
		}

	} else {
		//	bufChan = t.pongChannel
		prevPingPong = "PONG"
		tm = t.mpong.GetMatrics()
		if t.debugFlag {
			fmt.Println("work10secTimeout - reading from PONG BUFFER")
		}

	}

	for !exitLoop {

		select {
		case traceDataSample = <-*t.previousChan:
			//case traceData = <-bufChan:
			{
				if t.debugFlag {
					fmt.Println("work10SecTimeout t.calc10SecsLogEntry ChannelType: ", prevPingPong)
				}
				//t.calc10SecsLogEntry(traceData, channelPreviousActiveType2)
				t.calc10SecsLogEntry(traceDataSample, t.channelPreviousActiveType)

			}
		default:
			{
				//at this time the channel buffer is empty
				if t.debugFlag {
					fmt.Println("work10SecTimeout - Default - NO MORE DATA IN BUFFER")
				}

				var tMsg LoggerMsg

				valueMap := map[string]time.Duration{
					"AvgRoundTripResponseDuration":  tm.AvgRoundTripResponseDuration,
					"PeakRoundTripResponseDuration": tm.PeakRoundTripResponseDuration,
					"sumRoundTripResponseDuration":  tm.sumRoundTripResponseDuration}

				tMsg = LoggerMsg{logType: LOG_TYPE_10SEC_MSG,
					Marker: "NEW", Devider: true, LineFeed: false, numClients: t.numClients, pingPong: prevPingPong,
					totalMsgPer10Sec: tm.totalMsgsPer10Secs, avgContentSizeBytes: tm.AvgContentSizeBytes, Attributes: valueMap}
				t.logger.Log(tMsg)

				tm.Reset()

				exitLoop = true

			}
		}

	}
	//etime := time.Now()
	//extime := etime.Sub(stime)
	if t.debugFlag {
		var tMsg LoggerMsg
		tMsg = LoggerMsg{
			logType:         LOG_TYPE_MAX_MSG,
			maxNumberOfMsgs: t.totalMessageProcessDuringTest,
		}
		t.logger.Log(tMsg)
		fmt.Println("work10SecTimeout - calling work10Lock UNLOCK")
		fmt.Println("work10SecTimeout ExitLoop")
	}

	t.work10Lock.Unlock()

}

func (t *TraceEngine) run() {
	//fmt.Println("run")
	go func() {
		fmt.Println("Trace Engine Running")
		t.startedFlag = true
		t.tickerMgr.StartTicking()
		var traceData TraceDataSample

		for t.dontExit {

			select {
			case traceData = <-*t.activeChannel:
				{
					if t.debugFlag {
						fmt.Println("RUN - calling calc10SecsLogEntry", t.channelActiveType)
					}
					t.calc10SecsLogEntry(traceData, t.channelActiveType)

				}
			case t.dontExit = <-t.dontExitChannel:
				{
					fmt.Println("run t.dontExit = <- t.dontExitChannel")
				}
			case <-t.tickerMgr.TenSecsTicker.C:
				{
					//fmt.Println("run calling waitOnTraceEntryLock LOCK")

					t.waitOnTraceEntryLock.Lock()

					//fmt.Println("RUN - TIMEOUT RECEIVED")

					if t.tickerMgr.counter10SecTicks < t.tickerMgr.max10SecTickerInterval {

						t.tickerMgr.counter10SecTicks++

						if t.getChannelActive() == CHANNEL_TYPE_PING {
							if t.debugFlag {
								fmt.Println("Writing TIMEOUT : PING SWITCHING TO PONG")
							}

							t.setChannelActive(CHANNEL_TYPE_PONG)
							*t.activeChannel = t.pongChannel
							*t.previousChan = t.pingChannel
							t.channelPreviousActiveType = CHANNEL_TYPE_PING
						} else {
							if t.debugFlag {
								fmt.Println("Writing TIMEOUT : PONG SWITCHING TO PING")
							}
							t.setChannelActive(CHANNEL_TYPE_PING)
							*t.activeChannel = t.pingChannel
							*t.previousChan = t.pongChannel
							t.channelPreviousActiveType = CHANNEL_TYPE_PONG

						}
						//fmt.Println("run calling waitOnTraceEntryLock UNLOCK")
						t.waitOnTraceEntryLock.Unlock()
						fmt.Println("10 SECS TIMEOUT PERIOD COMPLETED")
						t.work10SecTimeout()

					} else {
						//t.logger.LogTestRunTotals(t.totalMessageProcessDuringTest)
						t.DoExit = true
						t.waitOnTraceEntryLock.Unlock()

						fmt.Println("Trace ready to exit out... ")
						time.Sleep(2 * time.Second)
						t.testDone <- true
					}

				}

			}

		}
	}()
}
