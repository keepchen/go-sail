package email

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Envelope email envelope
type Envelope struct {
	From     string                       //Sender email address
	To       string                       //Recipient email address
	Subject  string                       //Subject
	MimeType string                       //Email MIME type
	Body     string                       //Email content
	Callback func(e *Envelope, err error) //Callback function
}

// Sender email sender
type Sender interface {
	// Mount adds to the send queue
	Mount(index int, envelope *Envelope)
	// Emit starts the queue
	Emit()
	// Done marks sending as complete
	Done()
	// Initialize sender goroutines
	makeWorkers()
	// Perform the send operation
	send(ep *Envelope) error
}

// sendPool send pool
type sendPool struct {
	workers      []chan *Envelope
	throttle     time.Duration
	workersCount int
	conf         Conf
	wg           *sync.WaitGroup
	exit         chan struct{}
}

var _ Sender = &sendPool{}

// NewPool instantiates a pool
func NewPool(conf Conf) Sender {
	pool := &sendPool{
		conf:         conf,
		workersCount: conf.Workers,
		throttle:     time.Duration(conf.WorkerThrottleSeconds) * time.Second,
		wg:           &sync.WaitGroup{},
		exit:         make(chan struct{}),
	}
	pool.makeWorkers()

	return pool
}

// Mount adds to the send queue
func (p *sendPool) Mount(index int, envelope *Envelope) {
	p.wg.Add(1)
	go func(ep *Envelope) {
		p.workers[index%p.workersCount] <- ep
	}(envelope)
}

// Emit starts the queue
func (p *sendPool) Emit() {
	handler := func(index int, wk chan *Envelope) {
	LOOP:
		for {
			select {
			case ep := <-wk:
				if ep != nil && ep != (*Envelope)(nil) {
					err := p.send(ep)
					if err != nil {
						fmt.Println("[GO-SAIL] <Email> send failure via worker, error:", err)
					}
					//Handle the caller's callback
					if ep.Callback != nil {
						ep.Callback(ep, err)
					}
					p.wg.Done()

					if p.throttle > 0 {
						time.Sleep(p.throttle)
					}
				}
			case <-p.exit:
				break LOOP
			}
		}
	}
	for index := range p.workers {
		go handler(index, p.workers[index])
	}
}

// Done marks sending as complete
func (p *sendPool) Done() {
	p.wg.Wait()
	for index := range p.workers {
		close(p.workers[index])
	}
	p.exit <- struct{}{}
	close(p.exit)
}

// Initialize sender goroutines
func (p *sendPool) makeWorkers() {
	if p.workersCount < 1 {
		p.workersCount = 1
	}

	fmt.Printf("[GO-SAIL] <Email> Init {%d} worker(s), throttle {%d} second(s)", p.workersCount, p.conf.WorkerThrottleSeconds)

	var workers = make([]chan *Envelope, p.workersCount)
	for index := range workers {
		workers[index] = make(chan *Envelope)
	}

	p.workers = workers
}

// Perform the send operation
func (p *sendPool) send(ep *Envelope) error {
	if ep == nil || ep == (*Envelope)(nil) {
		return errors.New("ep is nil")
	}
	headers := map[string]string{
		"From":         ep.From,
		"Subject":      ep.Subject,
		"Content-Type": fmt.Sprintf("%s; charset=UTF-8", ep.MimeType),
		"To":           ep.To,
	}
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + ep.Body

	err := New(p.conf).SendMailWithTLS([]string{ep.To}, []byte(message))

	return err
}
