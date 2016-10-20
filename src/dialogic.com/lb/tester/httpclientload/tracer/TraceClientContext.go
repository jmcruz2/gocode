package tracer

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"time"

	"strconv"

	"github.com/tcnksm/go-httpstat"
	"strings"
)

type TraceClientContext struct {
	ClientContextId int
	MsgId           int
	traceLogger     TraceLogger
	cb              clientBucket
	tEngine         TraceEngine
	debugFlag       bool
}

func (t *TraceClientContext) Start(te TraceEngine, testParms CommandLineParams) {

	go func() {
		if t.debugFlag {
			fmt.Println("Starting Http Client Load Tester")
		}

		t.tEngine = te

		t.traceLogger = t.tEngine.logger

		t.cb = clientBucket{}
		//req, err := http.NewRequest("GET", *testParms.LbUri, nil)
		//
		//if err != nil {
		//	fmt.Println("PANIC CLIENTID: ", t.ClientContextId)
		//	panic(err)
		//}
		//
		//var statResult httpstat.Result
		//ctx := httpstat.WithHTTPStat(req.Context(), &statResult)
		//req = req.WithContext(ctx)

		//client := &http.Client{Transport: &t.cb}
		client := &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 7,
			},
		}

		trace := &httptrace.ClientTrace{

			ConnectStart: func(network, addr string) {
				//fmt.Println("Dial start")
				t.cb.connectStartTime = time.Now()
				t.cb.connectionStartDuration = t.cb.connectStartTime.Sub(t.cb.connectStartTime)

			},
			ConnectDone: func(network, addr string, err error) {
				//fmt.Println("ConnectDone")
				t.cb.connectDoneTime = time.Now()
				t.cb.connectDoneDuration = t.cb.connectDoneTime.Sub(t.cb.startingRequestTime)
			},
			GotConn: t.cb.GotConn,

			GotFirstResponseByte: func() {
				t.cb.calcTotalServerProcessing()
				t.cb.calcRcvLatency()

			},
			WroteHeaders: func() {
				//fmt.Println("Wrote headers")
			},
			WroteRequest: func(wr httptrace.WroteRequestInfo) {
				t.cb.wroteRequestTime = time.Now()
				t.cb.wroteRequestDuration = t.cb.wroteRequestTime.Sub(t.cb.gotConnectionTime)
				t.cb.calcXmitLatency()
			},
		}

		//create request here only once gets best output rate
		req, err := http.NewRequest("GET", *testParms.LbUri, nil)

		if err != nil {
			fmt.Println("PANIC CLIENTID: ", t.ClientContextId)
			panic(err)
		}

		var statResult httpstat.Result
		ctx := httpstat.WithHTTPStat(req.Context(), &statResult)
		req = req.WithContext(ctx)
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

		//reqClientId := strconv.Itoa(t.ClientContextId)
		//reqMsgId := strconv.Itoa(t.MsgId)
		//reqClientMsgIds := fmt.Sprintf("id=%s  no=%s", reqClientId, reqMsgId)
		//req.Header.Add("id", reqClientMsgIds)

		//for i := 0; i < *testParms.CallsPerSecond; i++ {
		t.debugFlag = *testParms.DbgFlag
		for {
			if !t.tEngine.DoExit {
				t.MsgId++
				t.cb.startingRequestTime = time.Now()

				if t.debugFlag {
					fmt.Println("Start client.DO HTTP Request")
				}

				//create request here slows down the number of messages sent by a large amount
				//req, err := http.NewRequest("GET", *testParms.LbUri, nil)
				//
				//if err != nil {
				//	fmt.Println("PANIC CLIENTID: ", t.ClientContextId)
				//	panic(err)
				//}
				//
				//var statResult httpstat.Result
				//ctx := httpstat.WithHTTPStat(req.Context(), &statResult)
				//req = req.WithContext(ctx)
				//req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
				//

				var reqClientId  string
				var reqMsgId 	 string
				var reqClientMsgIds string

				if t.MsgId == 1 {
					reqClientId = strconv.Itoa(t.ClientContextId)
					reqMsgId = strconv.Itoa(t.MsgId)
					reqClientMsgIds = fmt.Sprintf("id=%s  no=%s", reqClientId, reqMsgId)
					req.Header.Add("id", reqClientMsgIds)
				} else {
					req.Header.Del("id")
					reqClientId := strconv.Itoa(t.ClientContextId)
					reqMsgId := strconv.Itoa(t.MsgId)
					reqClientMsgIds := fmt.Sprintf("id=%s  no=%s", reqClientId, reqMsgId)
					req.Header.Add("id", reqClientMsgIds)
				}

				resp, err := client.Do(req)

				if err != nil {
					fmt.Println("MAJOR ERROR client.DO CLIENTID - error: ", t.ClientContextId, err)
					//panic(err)
				} else {
					respClientMsgIds := resp.Header.Get("id")

					if t.debugFlag {
						fmt.Println("Got back: ", respClientMsgIds)
					}

					respHeaderMsgIds := resp.Header.Get("id")
					if strings.Compare(respClientMsgIds, respHeaderMsgIds) != 0 {
						fmt.Println("Expected ClientId MsgId ", reqClientMsgIds, " Instead got back: ", respHeaderMsgIds)
					}

					io.Copy(ioutil.Discard, resp.Body)

					//io.Copy(os.Stdout, resp.Body)
					resp.Body.Close()
					if t.debugFlag {
						fmt.Println("Completed client.DO HTTP Request CLIENTID: ", t.ClientContextId, t.MsgId)
					}

					traceDataEntry := TraceDataEntry{
						newConnection:     t.cb.reusingConnection,
						clientID:          t.ClientContextId,
						msgID:             t.MsgId,
						latencyDuration:   t.cb.getTotalLatency(),
						serverDuration:    t.cb.getTotalServerProcessing(),
						totalRespDuration: t.cb.calcTotalResponseTime(),
					}

					if t.debugFlag {
						fmt.Println("calling Engine.TraceEntry ", t.ClientContextId, t.MsgId )
					}
					t.tEngine.TraceEntry(traceDataEntry)
				}
			} else {
				fmt.Println("Exiting ClientContext CLIENTID: ", t.ClientContextId)
				break
			}

		}

		//io.Copy(os.Stdout , resp.Body.Read())

		fmt.Println("TraceClientContext Done!")
	}()

}
