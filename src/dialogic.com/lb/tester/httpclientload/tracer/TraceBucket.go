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
	current           *http.Request
	ReusingConnection bool

	ConnectStartTime        time.Time
	ConnectionStartDuration time.Duration

	ConnectDoneTime     time.Time
	ConnectDoneDuration time.Duration
	GotConnectionTime   time.Time

	FirstByteRecievedTime              time.Time
	FirstByteRoundTripResponseDuration time.Duration

	WroteRequestTime     time.Time
	WroteRequestDuration time.Duration

	RoundTripResponseTime     time.Time
	RoundTripResponseDuration time.Duration
	ContentSizeBytes          int64
}

func (cObj *clientBucket) reset() {
	cObj.WroteRequestDuration = 0
	cObj.RoundTripResponseDuration = 0
	cObj.ContentSizeBytes = 0
}

/****
func (t *clientBucket) RoundTrip(req *http.Request) (*http.Response, error) {
	//fmt.Println("RoundTrip......")
	t.current = req
	resp, err := http.DefaultTransport.RoundTrip(req)
	return resp, err
}
***/

// GotConn prints whether the connection has been used previously
// for the current request.
func (t *clientBucket) GotConn(info httptrace.GotConnInfo) {
	//fmt.Println("Connection reused: ", info.Reused)
	t.ReusingConnection = info.Reused
	if t.ReusingConnection {
		t.ConnectDoneDuration = 0
		t.GotConnectionTime = time.Now()
	}
}
