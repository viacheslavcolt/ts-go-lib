package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type Listener interface {
	HandleEv(ev *Event)
}

type GroupConfig struct {
	Ln      Listener
	GroupId string
	Topic   string
	Workers int
	Timeout time.Duration
}

type group struct {
	cncls []context.CancelFunc
	wg    *sync.WaitGroup
}

func newGroup() *group {
	return &group{
		cncls: make([]context.CancelFunc, 0),
		wg:    &sync.WaitGroup{},
	}
}

func (g *group) makeCtx() context.Context {
	var (
		ctx  context.Context
		cncl context.CancelFunc
	)

	ctx, cncl = context.WithCancel(context.Background())

	g.cncls = append(g.cncls, cncl)

	return ctx
}

func (g *group) shutdown() {
	for _, cncl := range g.cncls {
		cncl()
	}
}

type ConsumerPool struct {
	brokers []string
	logger  zerolog.Logger

	groups []*group
}

func NewPool(brokers []string, logger zerolog.Logger) *ConsumerPool {
	return &ConsumerPool{
		brokers: brokers,
		logger:  logger,
	}
}

func (p *ConsumerPool) workerFn(ctx context.Context, r *kafka.Reader, l Listener, timeout time.Duration) {
	var (
		msg kafka.Message
		ev  Event
		// ok       bool
		// kafkaErr kafka.Error

		isRunning bool

		err error
	)

	isRunning = true

	for isRunning {
		select {
		case <-ctx.Done():
			isRunning = false
			break
		default:
			if msg, err = r.ReadMessage(ctx); err != nil {
				p.logger.Err(err).Msg("consumer error")
				isRunning = false
				continue
				// if kafkaErr, ok = err.(kafka.Error); ok {

				// }
			}

			if err = proto.Unmarshal(msg.Value, &ev); err != nil {
				p.logger.Err(err).Msg("unmarshal message error")
				continue
			}

			l.HandleEv(&ev)
		}
	}

	r.Close()
}

func (p *ConsumerPool) startGroup(g *group, cfg *GroupConfig) error {
	var (
		r *kafka.Reader
	)

	for i := 0; i < cfg.Workers; i++ {
		g.wg.Add(1)

		r = kafka.NewReader(kafka.ReaderConfig{
			Brokers: p.brokers,
			Topic:   cfg.Topic,
			GroupID: cfg.GroupId,
		})

		go p.runWorker(g, r, cfg)
	}

	return nil
}

func (p *ConsumerPool) runWorker(g *group, r *kafka.Reader, cfg *GroupConfig) error {
	p.workerFn(g.makeCtx(), r, cfg.Ln, cfg.Timeout)

	g.wg.Done()

	return nil
}

func (p *ConsumerPool) Start(cfg *GroupConfig) (*sync.WaitGroup, error) {
	var (
		g   *group
		err error
	)

	g = newGroup()

	if err = p.startGroup(g, cfg); err != nil {
		g.shutdown()
		return nil, err
	}

	p.groups = append(p.groups, g)

	return g.wg, nil
}

func (p *ConsumerPool) StopAll() {
	for _, g := range p.groups {
		g.shutdown()
	}
}
