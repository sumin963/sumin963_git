package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/sumin963/broadcast"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Global Declarations
	const QWheel_size int = 12  //const int QWheel_size=12;
	const QBumper_size int = 12 //const int QBumper_size=12;
	const QCliff_size int = 12  //const int QCliff_size=12;

	const spinRate int = 1 //const int spinRate=1;
	const TimeOut int = 1  //const int TimeOut=1;
	const CBmin int = 4    //const int CBmin=4
	const CBmax int = 5    //const int CBmax=5

	type sensor_Type int //typedef int[0,2] sensor_Type;
	type sensor_Pos int  //typedef int[0,2] sensor_Pos;

	type subscriber_QID int //typedef int[0,2] subscriber_QID;

	sensorMsg := make([][]chan interface{}, 3) //chan sensorMsg[sensor_Type][sensor_Pos];
	for i := range sensorMsg {
		sensorMsg[i] = make([]chan interface{}, 3)
	}
	//Go := make(chan interface{})
	//getMsgs := make(chan interface{})
	//spinOnce := make(chan interface{})  //
	var getMsgs broadcast.Broadcaster  //urgent broadcast chan getMsgs;
	var Go broadcast.Broadcaster       // urgent broadcast chan go;
	var spinOnce broadcast.Broadcaster //urgent broadcast chan spinOnce;

	var CBavali int = 0 //int[0,50] CBavail=0;

	//E := 0.01sec
	E := time.Millisecond * 100

	// Template sensor
	sensor := func(SensorType sensor_Type, SensorPos sensor_Pos) {
		sensorMsg[SensorType][SensorPos] = make(chan interface{})

		// sensor Declarations
		t := time.Now()     //clock t;
		ct := time.Since(t) // Cumulative clock t

	Init: // <Init Location & Edge>
		fmt.Printf("sensor(%v %v) Init\n", SensorType, SensorPos)
		ct = time.Since(t)
		<-time.After(time.Second - ct - E)             // Guard: t >= 1
		sensorMsg[SensorType][SensorPos] <- struct{}{} // Sync: sensorMsg[SensorType][SensorPos]!
		t = time.Now()                                 // Update: t = 0
		goto Init                                      // [Init --> Init]
	}

	// Template subscriberQueue
	subscriberQueue := func(SubscriberQID sensor_Type, SIZE int) {
		getMsgsls := getMsgs.Listen()
		var st sensor_Type = SubscriberQID // Guard: SubscriberQID==st and
		var pos sensor_Pos
		var MsgCount int = 0 //int[0, SIZE] MsgCount=0;

	Init: // <Init Location & Edge>
		fmt.Printf("subqueue(%v) Init\n", SubscriberQID)
		select {
		case <-sensorMsg[st][pos]: // Sync: sensorMsg[st][pos]?
			switch {
			case SubscriberQID == st && MsgCount >= SIZE: // Guard: SubscriberQID==st and MsgCount>=SIZE
				goto Overflow // [Init --> Overflow]
			case SubscriberQID == st && MsgCount < SIZE: // Guard: SubscriberQID==st and MsgCount<SIZE
				MsgCount++ // Update: MsgCount++
				CBavali++  // Update: CBavali++
				goto Ok    // [Init --> Ok]
			}
		case <-getMsgsls.Ch: // Sync: getMsgs?
			fmt.Printf("subq(%v) gm1\n", SubscriberQID)
			MsgCount = 0 // Update: MsgCount=0
			goto Init    // [Init --> Init]
		}

	Ok: // <Init Location & Edge>
		fmt.Printf("subqueue(%v) Ok\n", SubscriberQID)
		Go.Send(struct{}{}) // Sync: go!
		goto Init           // [Ok --> Init]

	Overflow: // <Init Location & Edge>
		fmt.Printf("subqueue(%v) Overflow\n", SubscriberQID)
		select {
		case <-getMsgsls.Ch: // Sync: getMsgs?
			fmt.Printf("subq(%v ) gm2\n", SubscriberQID)
			MsgCount = 0 // Update: x = 0
			goto Init    // [Overflow --> Init]
		case <-sensorMsg[st][pos]: // Sync: sensorMsg[st][pos]?
			goto Overflow // [Overflow --> Overflow]
		}
	}

	// Template safetyContrillrtUpdate
	safetyContrillrtUpdate := func() {
		t := time.Now()     // clock t
		ct := time.Since(t) // Cumulative clock t

	Init: // <Init Location & Edge>
		fmt.Printf("safetyContrillrtUpdate Init\n")
		ct = time.Since(t)
		<-time.After(time.Duration(spinRate)*time.Second - ct - E) // Guard: x >= 1
		select {
		case <-time.After(0):
			spinOnce.Send(struct{}{}) // Sync: spinOnce!
			goto spinLoc              // [Init --> spinLoc]
		case <-time.After(E):
			goto Alarm // Invariant Violation
		}
	spinLoc: // <spinLoc Location & Edge>
		fmt.Printf("safetyContrillrtUpdate spinLoc\n")
		t = time.Now() // Update: t = 0
		goto Init      // [spinLoc --> Init]

	Alarm: // <Alarm>
		fmt.Println("Invariant Violation!sc")
	}

	// Template callBackQueue
	callBackQueue := func() {
		// callBackQueue Declarations
		t1 := time.Now()      // clock t1
		ct1 := time.Since(t1) // Cumulative clock cx1
		t2 := time.Now()      // clock t2
		ct2 := time.Since(t2) // Cumulative clock ct2

		var MinCBTime int = 0 //int[0,10] MinCBTime;
		var MaxCBTime int = 0 //int[0,10] MaxCBTime;

		setCBTime := func() { //void setCBTime()
			MaxCBTime = CBmax
			MinCBTime = CBmin
		}
		Gols := Go.Listen()
		spinOncels := spinOnce.Listen()

		CBavalich := make(chan interface{})

		CBavaliguard := func() {
			for {
				if CBavali >= 1 {
					close(CBavalich)
					break
				}
			}
		}

	Init: // <Init Location & Edge>
		fmt.Printf("CBqueue Init\n")
		<-spinOncels.Ch  // Sync: spinOnce?
		goto anonnymous1 // [Init --> annonymous1]

	anonnymous1: // <annonymous1 Location & Edge>
		fmt.Printf("CBqueue anm1\n")
		switch {
		case CBavali == 0: // Guard: CBavali==0
			t2 = time.Now() // Update: t2 = 0
			goto wait       // [annonymous1 --> wait]
		case CBavali > 0: // Guard: CBavali>0
			getMsgs.Send(struct{}{}) // Sync: getMsgs!
			setCBTime()              // Update: setCBTime()
			t1 = time.Now()          // Update: t1 = 0
			CBavali = 0              // Update: CBavali =0
			goto CBprocess           // [annonymous1 --> wait]
		}

	CBprocess: // <CBprocess Location & Edge>
		fmt.Printf("CBqueue CBprocess\n")
		ct1 = time.Since(t1)
		<-time.After(time.Duration(MinCBTime)*time.Second - ct1 - E) // Guard: t1>=MinCBTime
		ct1 = time.Since(t1)
		select {
		case <-time.After(0):
			goto Init // [CBprocess --> Init]
		case <-time.After(time.Duration(MaxCBTime)*time.Second - ct1):
			goto Alarm // Invariant Violation
		}

	wait: // <wait Location & Edge>
		fmt.Printf("CBqueue wait\n")
		ct2 = time.Since(t2)
		CBavalich = make(chan interface{})
		CBavaliguard()
		select {
		case <-time.After(time.Duration(TimeOut)*time.Second - ct2 - E): // Guard: t2==TimeOut
			goto Init // [wait --> Init]
		case <-CBavalich: // Guard: CBavail>=1
			<-Gols.Ch        // Sync: go?
			goto anonnymous2 // [wait --> anonnymous2]
		case <-time.After(time.Duration(TimeOut)*time.Second - ct2):
			goto Alarm // Invariant Violation
		}

	anonnymous2: // <anonnymous2 Location & Edge>
		fmt.Printf("CBqueue anm2\n")
		getMsgs.Send(struct{}{}) // Sync: getMsgs!
		setCBTime()              // Update: setCBTime()
		t1 = time.Now()          // Update: t1 = 0
		CBavali = 0              // Update: CBavali = 0
		goto CBprocess           // [anonnymous2 --> CBprocess]

	Alarm: // <Alarm>
		fmt.Println("Invariant Violation!cq")
	}

	// System declarations
	go sensor(0, 0)                    // Wheel_Left=sensor(0,0);
	go sensor(0, 1)                    // Wheel_Right=sensor(0,1);
	go subscriberQueue(0, QWheel_size) // QWheel=subscriberQueue(0,QWheel_size);

	go sensor(1, 0)                     // Bumper_Left=sensor(1,0);
	go sensor(1, 1)                     // Bumper_Center=sensor(1,1);
	go sensor(1, 2)                     // Bumper_Right=sensor(1,2);
	go subscriberQueue(1, QBumper_size) // QBumper=subscriberQueue(1,QBumper_size);

	go sensor(2, 0)                    // Cliff_Left=sensor(2,0);
	go sensor(2, 1)                    // Cliff_Center=sensor(2,1);
	go sensor(2, 2)                    // Cliff_Right=sensor(2,2);
	go subscriberQueue(2, QCliff_size) // QCliff=subscriberQueue(2,QCliff_size);

	go safetyContrillrtUpdate() // SafetyControllerUpdate=safetyControllerUpdate();
	go callBackQueue()          // CBQueue=callBackQueue();

	<-time.After(time.Second * 600)
}
