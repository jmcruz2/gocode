package main

import (
	"flag"
	"fmt"
	"os"

	"dialogic.com/lb/tester/httpclientload/tracer"
)

var gclientNextID = 0

var traceEngine tracer.TraceEngine
var testDone chan bool

var CmdLineParams tracer.CommandLineParams

func main() {
	fmt.Printf("%s\n", "Entering main - Http Client Load Tester")

	if len(os.Args) < 9 {
		fmt.Fprintf(os.Stderr, "Usage: -p [IP:port] -n -r -d\n")
		os.Exit(1)
	}

	CmdLineParams = tracer.CommandLineParams{}

	CmdLineParams.NumOfClients = flag.Int("n", 10, "Number of Clients or client threads")
	CmdLineParams.CallsPerSecond = flag.Int("r", 100, "Calls Per Second")
	CmdLineParams.RepetionInMin = flag.Int("d", 1, "Number of repetition in minutes")
	CmdLineParams.LbUri = flag.String("p", "localhost:8080", "LB IP Address and Port")
	CmdLineParams.DbgFlag = flag.Bool("debug", false, "Debug Mode")
	flag.Parse()
	//fmt.Println("OS var cnt: ", len(os.Args))
	fmt.Println("Number of Clients: ", *CmdLineParams.NumOfClients)
	fmt.Println("Calls per seconds: ", *CmdLineParams.CallsPerSecond)
	fmt.Println("Number of repetition in minutes: ", *CmdLineParams.RepetionInMin)
	fmt.Println("LB IP Addfess and Port: ", *CmdLineParams.LbUri)
	fmt.Println("Debug Mode: ", *CmdLineParams.DbgFlag)

	testDone = make(chan bool, 0)
	traceEngine = tracer.TraceEngine{}
	traceEngine.Ctor(100000, *CmdLineParams.RepetionInMin, *CmdLineParams.NumOfClients, testDone, *CmdLineParams.DbgFlag, *CmdLineParams.CallsPerSecond)

	traceEngine.Start()

	for n := 0; n < *CmdLineParams.NumOfClients; n++ {
		//fmt.Println("calling go startTest()")
		startTest(CmdLineParams, n+1)
	}

	<-testDone
	fmt.Println("Load Test Completed!")

}

func startTest(cmdLineParameters tracer.CommandLineParams, clientID int) {

	//fmt.Println("Starting Http Client Load Tester CLIENTID: ", clientID)
	gclientNextID = clientID
	tcContext := tracer.TraceClientContext{ClientContextId: clientID}
	//fmt.Println("---------Starting request!------------")
	go tcContext.Start(traceEngine, cmdLineParameters)
	// fmt.Println("Main startTest Done! CLIENTID: ", clientID)

}
