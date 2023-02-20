package email

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Envelope 信封
type Envelope struct {
	From     string                       //发送人邮箱
	To       string                       //收件人邮件
	Subject  string                       //主题
	MimeType string                       //邮件类型
	Body     string                       //邮件内容
	Callback func(e *Envelope, err error) //回调函数
}

// Pool 发送池
type Pool struct {
	workers  []chan *Envelope
	throttle time.Duration
	count    int
	conf     Conf
	wg       *sync.WaitGroup
}

// NewPool 实例化
func NewPool(conf Conf) *Pool {
	pool := &Pool{
		conf:     conf,
		count:    conf.Workers,
		throttle: time.Duration(conf.WorkerThrottleSeconds) * time.Second,
		wg:       &sync.WaitGroup{},
	}
	pool.makeWorkers()

	return pool
}

// Mount 加入发送队列
func (p *Pool) Mount(index int, envelope *Envelope) {
	p.wg.Add(1)
	go func(ep *Envelope) {
		p.workers[index%p.count] <- ep
	}(envelope)
}

// Emit 启动队列
func (p *Pool) Emit() {
	for index, worker := range p.workers {
		go func(index int, wk chan *Envelope) {
			for {
				select {
				case ep := <-wk:
					err := p.send(ep)
					if err != nil {
						log.Println("[!] send failure via worker, error:", err)
					}
					//处理上层回调函数
					if ep.Callback != nil {
						ep.Callback(ep, err)
					}
					p.wg.Done()
					if p.throttle > 0 {
						time.Sleep(p.throttle)
					}
				}
			}
		}(index, worker)
	}
}

// Done 发送完成
func (p *Pool) Done() {
	p.wg.Wait()
	for index, _ := range p.workers {
		close(p.workers[index])
	}
}

// 初始化发送协程
func (p *Pool) makeWorkers() {
	if p.count < 1 {
		p.count = 1
	}

	log.Printf("[!] Init {%d} worker(s), throttle {%d} second(s)", p.count, p.conf.WorkerThrottleSeconds)

	var workers = make([]chan *Envelope, p.count)
	for index, _ := range workers {
		workers[index] = make(chan *Envelope)
	}

	p.workers = workers
}

// 执行发送操作
func (p *Pool) send(ep *Envelope) error {
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
