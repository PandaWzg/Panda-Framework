package pool

import (
	"sync"
)

type Pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

// 创建并发控制池, 设置并发数量与总数量
func New(cap int) *Pool {
	if cap < 1 {
		cap = 1
	}
	p := &Pool{
		queue: make(chan int, cap),
		wg:    new(sync.WaitGroup),
	}
	return p
}
func (p *Pool) AddCount(total int) {
	p.wg.Add(total)
}

// 向并发队列中添加一个
func (p *Pool) AddOne() {
	p.queue <- 1
	p.wg.Add(1)
}

// 并发队列中释放一个, 并从总数量中减去一个
func (p *Pool) DelOne() {
	<-p.queue
	p.wg.Done()
}

func (p *Pool) Done() {
	p.wg.Done()
}

func (p *Pool) Wait() {
	//close(p.queue)
	p.wg.Wait()
}
