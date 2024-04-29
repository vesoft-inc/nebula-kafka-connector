package async

import (
	"context"
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type Group struct {
	ctx         context.Context
	concurrency int
	workers     []func(stopCh chan interface{})
	event       string
}

func NewGroup(ctx context.Context, concurrency int, event string) *Group {
	return &Group{
		ctx:         ctx,
		concurrency: concurrency,
		event:       event,
	}
}

func (g *Group) Add(worker func(stopCh chan interface{})) {
	g.workers = append(g.workers, worker)
}

func (g *Group) Wait() error {
	if g.workers == nil || len(g.workers) == 0 {
		return nil
	}

	stopCh := make(chan interface{}, len(g.workers))
	concurCh := make(chan interface{}, g.concurrency)

	go func() {
		for _, worker := range g.workers {
			concurCh <- struct{}{}
			go worker(stopCh)
		}
	}()

	var errs []error
	res := 0
	for {
		select {
		case <-g.ctx.Done():
			// TODO: notify ongoing worker to stop
			return fmt.Errorf("group waiting time out")
		case val := <-stopCh:
			switch v := val.(type) {
			case error:
				errs = append(errs, v)
			case []error:
				errs = append(errs, v...)
			default:
			}
			<-concurCh
			res++
			log.Infof("%s finished, progress: %d/%d", g.event, res, len(g.workers))
			if res == len(g.workers) {
				if len(errs) == 0 {
					return nil
				}
				return utils.ErrorAggregate(errs)
			}
		}
	}
}
