package tracer

import (
	"fmt"
	"time"
)

type TraceTickerMgr struct {
	max10SecTickerInterval int
	counter10SecTicks      int
	TenSecsTicker          *time.Ticker
}

func (tm *TraceTickerMgr) Ctor(t int) {

	tm.max10SecTickerInterval = t * 6 //say t =10 min total of ticker interval  there are 6 ten seconds periods in a minute
	fmt.Println("TraceTicker number of repetions: ", t, " Number of 10 sec periods: ", tm.max10SecTickerInterval)
	tm.counter10SecTicks = 0

}

func (tm *TraceTickerMgr) StartTicking() {
	//fmt.Println("Starting 10 sec Ticker")
	tm.TenSecsTicker = time.NewTicker(time.Second * 10)
}

func (tm *TraceTickerMgr) StopTicking() {
	tm.TenSecsTicker.Stop()
}

func (tm *TraceTickerMgr) GetMax10SecTickerInterval() int {
	return tm.max10SecTickerInterval
}

func (tm *TraceTickerMgr) GetCounter10SecTicks() int {
	return tm.counter10SecTicks
}

func (tm *TraceTickerMgr) SetCounter10SecTicks(val int) int {
	tm.counter10SecTicks = val
	return tm.counter10SecTicks
}

func (tm *TraceTickerMgr) areWeDone() bool {
	weDone := false

	if tm.counter10SecTicks >= tm.max10SecTickerInterval {
		weDone = true
	}
	return weDone
}
