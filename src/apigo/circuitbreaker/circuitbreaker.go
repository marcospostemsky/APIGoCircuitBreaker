package circuitbreaker

import (
	"time"
	"fmt"
	"net/http"
	"../utils"
)


const (
	STATE_CLOSE         State        = 1
	STATE_HALF_OPEN     State        = 2
	STATE_OPEN          State        = 3
	//SERVICE_AVAILABLE   StateService = 1
	//SERVICE_UNAVAILABLE StateService = 0
)

type State int

type CircuitBreaker struct {
	Name          string
	MaxRequests   int
	Timeout       time.Duration

	State      		State
	ErrCounts     	int

	ChainStatus 	chan string
	ChainTimeout	chan time.Time
	ChainCount		chan string
}


func NewCircuitBreaker(name string, maxrequest int, timeout time.Duration, errcounts int ) *CircuitBreaker {
	cb := &CircuitBreaker{
		Name:name,
		MaxRequests:maxrequest,
		Timeout:timeout,
		ErrCounts:errcounts,
		State:STATE_CLOSE,
		ChainStatus:make(chan string, 2),
		ChainCount:make(chan string,2),
	}
	return cb
}


func (cb *CircuitBreaker) SetState(){
	go func() {
		for{
			<- cb.ChainStatus
			if cb.State == STATE_HALF_OPEN {
				fmt.Println("State: HALF-OPEN")
				response, err := http.Get(utils.UrlPing)
				if err != nil {
					cb.State = STATE_OPEN
				}
				if  response != nil {
					if response.StatusCode == 200 {
						cb.State = STATE_CLOSE
					}
				}
			}
			if cb.State == STATE_OPEN {
				fmt.Println("State: OPEN")
				cb.ChainTimeout <- <- time.After(cb.Timeout)

			}
			if cb.State == STATE_CLOSE {
				fmt.Println("State: CLOSE")
				cb.Reset()
			}
		}
	}()

	go func() {
		for{
			msg := <- cb.ChainCount
			cb.Counter(msg)
		}
	}()


	go func() {
		cb.ChainTimeout = make(chan time.Time)
		for {
			cb.WaitTimeOut()
		}
	}()
}


func (cb *CircuitBreaker) Reset(){
	cb.ErrCounts = 0
}

func (cb *CircuitBreaker) Counter(msg string){
	if msg == "ERROR"{
		cb.ErrCounts ++
	}
	if msg == "OK"{
		cb.ErrCounts = 0
	}
	if cb.ErrCounts >= cb.MaxRequests{
		cb.State=STATE_OPEN
		cb.ChainStatus <- "Status"
	}
}

func (cb *CircuitBreaker) WaitTimeOut (){
	<-cb.ChainTimeout
	if cb.State == STATE_OPEN {
		cb.State = STATE_HALF_OPEN
		cb.ChainStatus <- "Status"
	}
}

func (s State) toString() string {
	switch s {
	case STATE_CLOSE:
		return "close"
	case STATE_HALF_OPEN:
		return "half_open"
	case STATE_OPEN:
		return "open"
	default:
		return "undefine"
	}
}

