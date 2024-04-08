package printer

import (
	"testing"

	"github.com/go-playground/assert/v2"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type dummyPlan struct{}

var _ nebula.PlanInfo = &dummyPlan{}

func (d *dummyPlan) Id() []byte {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) Name() []byte {
	return []byte("dummy")
}

func (d *dummyPlan) Details() []byte {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) TimeMs() float64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) Rows() int64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) MemoryKib() float64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) BlockedMs() float64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) QueuedMs() float64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) ConsumeMs() float64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) ProduceMs() float64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) FinishMs() float64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) Batches() int64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) Concurrency() int64 {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) OtherStatsJson() []byte {
	panic("not implemented") // TODO: Implement
}

func (d *dummyPlan) Children() []nebula.PlanInfo {
	panic("not implemented") // TODO: Implement
}

func TestGetField(t *testing.T) {
	p := &dummyPlan{}
	v := fieldValue(p, "Name")
	assert.Equal(t, []byte("dummy"), v)

}
