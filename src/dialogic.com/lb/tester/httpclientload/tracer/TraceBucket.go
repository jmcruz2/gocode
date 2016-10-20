package tracer

import (
	"net/http"
	"net/http/httptrace"
	"time"
)

//total response time = total latency + server processing time
//server processing = gotFirstResponseTime - startingRequestTime
//total response time =  (nowTime after [io.Copy(ioutil.Discard, resp.Body) and resp.Body.Close()]) - startingRequestTime
//latency = total response time - server processing time

type clientBucket struct {
	current             *http.Request
	startingRequestTime time.Time
	reusingConnection   bool

	connectStartTime        time.Time
	connectionStartDuration time.Duration

	connectDoneTime     time.Time
	connectDoneDuration time.Duration
	gotConnectionTime	time.Time

	wroteRequestTime time.Time
	wroteRequestDuration	time.Duration

	totalServerProcessingDuration time.Duration
	totalResponseTimeDuration     time.Duration

	xmitLatencyDuration	      time.Duration
	rcvLatencyDuration	      time.Duration
	totalLatencyDuration          time.Duration

	totalServerProcessingTime time.Time
	totalResponseTime         time.Time
}

func (t *clientBucket) getTotalServerProcessing() time.Duration {
	return t.totalServerProcessingDuration
}

func (t *clientBucket) getTotalLatency() time.Duration {
	return t.totalLatencyDuration
}

func (t *clientBucket) getTotalResponseTime() time.Duration {
	return t.totalResponseTimeDuration
}

func (t *clientBucket) calcTotalServerProcessing() time.Duration {

	t.totalServerProcessingTime = time.Now()				//got first response byte
	t.totalServerProcessingDuration = t.totalServerProcessingTime.Sub(t.wroteRequestTime)
	return t.totalServerProcessingDuration
}


func (t *clientBucket) calcXmitLatency() time.Duration {

	if t.reusingConnection {
		t.xmitLatencyDuration = t.wroteRequestDuration
	}else {
		t.xmitLatencyDuration = t.wroteRequestTime.Sub(t.connectStartTime)
	}
	return t.xmitLatencyDuration
}

func (t *clientBucket) calcRcvLatency() time.Duration {
	endTime := time.Now()
	t.rcvLatencyDuration = endTime.Sub(t.totalServerProcessingTime)		//got first response byte
	t.totalLatencyDuration = t.xmitLatencyDuration + t.rcvLatencyDuration
	return t.rcvLatencyDuration
}

func (t *clientBucket) calcTotalResponseTime() time.Duration {
	t.totalResponseTimeDuration = t.totalLatencyDuration + t.totalServerProcessingDuration
	return t.totalResponseTimeDuration
}

func (t *clientBucket) calcTotalLatency() time.Duration {
	return t.totalLatencyDuration
}

func (t *clientBucket) RoundTrip(req *http.Request) (*http.Response, error) {
	//fmt.Println("RoundTrip......")
	t.current = req
	resp, err := http.DefaultTransport.RoundTrip(req)
	return resp, err
}

// GotConn prints whether the connection has been used previously
// for the current request.
func (t *clientBucket) GotConn(info httptrace.GotConnInfo) {
	//fmt.Println("Connection reused: ", info.Reused)
	t.reusingConnection = info.Reused
	if t.reusingConnection {
		t.connectDoneDuration 	= 0
		t.gotConnectionTime  = time.Now()

	}
}
