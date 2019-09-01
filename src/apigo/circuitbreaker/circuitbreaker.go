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
	ERROR				Count		 = 4
	OK					Count		 = 5
)

type State int
type Count int

type CircuitBreaker struct {
	Name          string
	MaxRequests   int
	Timeout       time.Duration

	State      		State
	ErrCounts     	int

	ChainStatus 	chan string
	ChainTimeout	chan time.Time
}


func NewCircuitBreaker(name string, maxrequest int, timeout time.Duration) *CircuitBreaker {
	cb := &CircuitBreaker{
		Name:name,
		MaxRequests:maxrequest,
		Timeout:timeout,
		ErrCounts:0,
		State:STATE_CLOSE,
		ChainStatus:make(chan string, 2),
	}
	return cb
}

// Función principal de Circuit  breaker, encargada de lanzar las go function que orquestan la apertura o no del interruptor.
func (cb *CircuitBreaker) SetState(){
	go func() {
		for{
			<- cb.ChainStatus
			if cb.State == STATE_HALF_OPEN {
				fmt.Println("State: HALF-OPEN")
				cb.ping()
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
		cb.ChainTimeout = make(chan time.Time)
		for {
			cb.WaitTimeOut()
		}
	}()
}


func (cb *CircuitBreaker) Reset(){
	cb.ErrCounts = 0
}


// Cuenta los errores secuenciales que se producen en el request. Se realizó en principio con una go func, pero se generaban
// problemas de sincronismo.
func (cb *CircuitBreaker) Counter(msg Count){
	if msg == ERROR {
		cb.ErrCounts ++
	}
	if msg == OK {
		cb.ErrCounts = 0
	}
	if cb.ErrCounts >= cb.MaxRequests{
		cb.State=STATE_OPEN
		cb.ChainStatus <- "Status"
	}
}


//Bloquea el canal hasta que llegue el time out, para cambiar el estado a Half-Open
func (cb *CircuitBreaker) WaitTimeOut (){
	<-cb.ChainTimeout
	if cb.State == STATE_OPEN {
		cb.State = STATE_HALF_OPEN
		cb.ChainStatus <- "Status"
	}
}

//Realiza el ping para verificar si se ha levantado la API, se implementa también un timeout utilizando el client de HTTP.

func (cb *CircuitBreaker) ping(){
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Get(utils.UrlPing)
	if err != nil {
		cb.State = STATE_OPEN
	}
	if  response != nil {
		if response.StatusCode == 200 {
			cb.State = STATE_CLOSE
		}
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

