package tracer

import "time"

type Total10SecTraceMatrics struct {

	/******** INSTANCE *****************/

	totalMsgsPer10Secs            int64
	AvgRoundTripResponseDuration  time.Duration
	PeakRoundTripResponseDuration time.Duration
	sumRoundTripResponseDuration  time.Duration
	AvgContentSizeBytes           int64
}

func (m *Total10SecTraceMatrics) GetMatrics() *Total10SecTraceMatrics {
	return m
}

func (m *Total10SecTraceMatrics) Reset() {

	//fmt.Println("Reset: TraceMatrics")
	m.totalMsgsPer10Secs = 0
	m.AvgRoundTripResponseDuration = 0
	m.AvgContentSizeBytes = 0
	m.PeakRoundTripResponseDuration = 0
}

type Imatrics interface {
	GetMatrics() *Total10SecTraceMatrics
	Reset()
}

type pingMatrics struct {
	Total10SecTraceMatrics
}

type pongMatrics struct {
	Total10SecTraceMatrics
}
