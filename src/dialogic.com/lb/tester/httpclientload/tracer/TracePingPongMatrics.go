package tracer

import "time"

type TraceMatrics struct {

	/******** INSTANCE *****************/
	minLatencyDuration time.Duration
	maxLatencyDuration time.Duration
	avgLatencyDuration time.Duration

	minServerDuration time.Duration
	maxServerDuration time.Duration
	avgServerDuration time.Duration

	minTotalRespDuration time.Duration
	maxTotalRespDuration time.Duration
	avgTotalRespDuration time.Duration

	/***********  TOTALS AVGS **************/
	avgLatencyTotal time.Duration
	avgLatencyCnt   int

	avgServerTotal time.Duration
	avgServerCnt   int

	avgTotalRespTotal time.Duration
	avgTotalRespCnt   int

	totalMsgsPer10Secs uint64
}

func (m *TraceMatrics) GetMatrics() *TraceMatrics {
	return m
}

func (m *TraceMatrics) Reset() {

	//fmt.Println("Reset: TraceMatrics")
	m.minLatencyDuration = 0
	m.maxLatencyDuration = 0
	m.avgLatencyDuration = 0
	m.minServerDuration = 0
	m.maxServerDuration = 0
	m.avgServerDuration = 0
	m.minTotalRespDuration = 0
	m.maxTotalRespDuration = 0
	m.avgTotalRespDuration = 0

	m.avgLatencyTotal = 0
	m.avgLatencyCnt = 0

	m.avgServerTotal = 0
	m.avgServerCnt = 0

	m.avgTotalRespTotal = 0
	m.avgTotalRespCnt = 0

	m.totalMsgsPer10Secs = 0
}

type Imatrics interface {
	GetMatrics() *TraceMatrics
	Reset()
}

type pingMatrics struct {
	TraceMatrics
}

type pongMatrics struct {
	TraceMatrics
}
