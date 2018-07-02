package bean

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/yarencheng/gospring/interfaces"
	"github.com/yarencheng/gospring/v1"
)

type Broadcaster struct {
	uuid     uuid.UUID
	id       string
	sourceID string
	size     int
	in       reflect.Value
	ctx      interfaces.ApplicationContextI
	stop     chan int
	wg       sync.WaitGroup
	// shared variables
	lock    sync.Mutex
	outs    []reflect.Value
	started bool
}

func V1BroadcastParser(ctx interfaces.ApplicationContextI, config interface{}) (interfaces.BeanI, error) {
	c, ok := config.(*v1.Broadcast)
	if !ok {
		return nil, fmt.Errorf("[%v] can not convert to [%v]", reflect.TypeOf(config), reflect.TypeOf(&v1.Broadcast{}))
	}

	b := &Broadcaster{
		uuid:     uuid.NewV4(),
		id:       c.ID,
		sourceID: c.SourceID,
		size:     c.Size,
		outs:     make([]reflect.Value, 0),
		ctx:      ctx,
		stop:     make(chan int),
		started:  false,
	}

	return b, nil
}

func (bc *Broadcaster) GetUUID() uuid.UUID {
	return bc.uuid
}

func (bc *Broadcaster) GetID() string {
	return bc.id
}

func (bc *Broadcaster) GetValue() (reflect.Value, error) {
	if !bc.in.IsValid() {
		if err := bc.initSourceChannel(); err != nil {
			return reflect.Value{}, fmt.Errorf("Initialize source channel failed, err: [%v]", err)
		}
	}

	inTypet := bc.in.Type().Elem()
	fmt.Printf("aaaaaa %v\n", inTypet)
	tvpe := reflect.ChanOf(reflect.BothDir, inTypet)
	out := reflect.MakeChan(tvpe, bc.size)

	bc.lock.Lock()
	defer bc.lock.Unlock()

	bc.outs = append(bc.outs, out)

	if !bc.started {
		bc.startBroadcast()
		bc.started = true
	}

	return out, nil
}

func (bc *Broadcaster) Stop(ctx context.Context) error {
	bc.stop <- 1

	wait := make(chan int, 1)
	go func() {
		bc.wg.Wait()
		wait <- 1
	}()

	select {
	case <-ctx.Done():
	case <-wait:
	}

	return ctx.Err()
}

func (bc *Broadcaster) initSourceChannel() error {
	bean, ok := bc.ctx.GetBeanByID(bc.sourceID)
	if !ok {
		return fmt.Errorf("Can't get bean with ID [%v]", bc.sourceID)
	}

	value, err := bean.GetValue()
	if err != nil {
		return fmt.Errorf("Can't get channel with ID [%v]", bc.sourceID)
	}

	bc.in = value

	return nil
}

func (bc *Broadcaster) startBroadcast() error {

	bc.wg.Add(1)
	go func() {
		defer bc.wg.Done()
		for {
			cases := make([]reflect.SelectCase, 2)
			cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: bc.in}
			cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(bc.stop)}
			chosen, value, ok := reflect.Select(cases)
			switch chosen {
			case 0:
				if !ok {
					break
				}
			case 1:
				break
			}

			bc.lock.Lock()
			defer bc.lock.Unlock()

			for _, out := range bc.outs {
				if out.Len() >= bc.size {
					// channel is full, skip this channel
					// TODO: refine
					continue
				}
				out.Send(value)
			}
		}

	}()

	return nil
}
