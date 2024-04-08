package nebula_ng

import "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/graph"

type (
	summary struct {
		summary *graph.Summary
	}
	planInfo struct {
		planInfo *graph.PlanInfo
	}
)

func (s *summary) BuildTimeUs() int64 {
	return s.summary.BuildTimeUs
}

func (s *summary) OptimizeTimeUs() int64 {
	return s.summary.OptimizeTimeUs
}

func (s *summary) LatencyUs() int64 {
	return s.summary.LatencyUs
}

func (s *summary) Preamble() []byte {
	return s.summary.Preamble
}

func (s *summary) PlanInfo() PlanInfo {
	return &planInfo{s.summary.GetPlanInfo()}
}

func (p *planInfo) Id() []byte {
	return p.planInfo.Id
}

func (p *planInfo) Name() []byte {
	return p.planInfo.Name
}

func (p *planInfo) Details() []byte {
	return p.planInfo.Details
}

func (p *planInfo) TimeMs() float64 {
	return p.planInfo.TimeMs
}

func (p *planInfo) Rows() int64 {
	return p.planInfo.Rows
}

func (p *planInfo) MemoryKib() float64 {
	return p.planInfo.MemoryKib
}

func (p *planInfo) BlockedMs() float64 {
	return p.planInfo.BlockedMs
}

func (p *planInfo) QueuedMs() float64 {
	return p.planInfo.QueuedMs
}

func (p *planInfo) ConsumeMs() float64 {
	return p.planInfo.ConsumeMs
}

func (p *planInfo) ProduceMs() float64 {
	return p.planInfo.ProduceMs
}

func (p *planInfo) FinishMs() float64 {
	return p.planInfo.FinishMs
}

func (p *planInfo) Batches() int64 {
	return p.planInfo.Batches
}

func (p *planInfo) Concurrency() int64 {
	return p.planInfo.Concurrency
}

func (p *planInfo) OtherStatsJson() []byte {
	return p.planInfo.OtherStatsJson
}

func (p *planInfo) Children() []PlanInfo {
	children := make([]PlanInfo, 0, len(p.planInfo.Children))
	for _, child := range p.planInfo.Children {
		children = append(children, &planInfo{planInfo: child})
	}
	return children
}
