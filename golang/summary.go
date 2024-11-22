package nebula_ng

import (
	"github.com/vesoft-inc/nebula-ng-tools/golang/internal/generated_code/v5.0.0/proto/graph"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

type (
	summary struct {
		summary *graph.Summary
	}
	planInfo struct {
		planInfo *graph.PlanInfo
	}
	queryStats struct {
		queryStats *graph.QueryStats
	}
)

func (s *summary) BuildTimeUs() int64 {
	return s.summary.ElapsedTime.BuildTimeUs
}

func (s *summary) OptimizeTimeUs() int64 {
	return s.summary.ElapsedTime.OptimizeTimeUs
}

func (s *summary) SerializeTimeUs() int64 {
	return s.summary.ElapsedTime.SerializeTimeUs
}

func (s *summary) TotalServerTimeUs() int64 {
	return s.summary.ElapsedTime.TotalServerTimeUs
}

func (s *summary) ExplainType() string {
	return string(s.summary.ExplainType)
}

func (s *summary) PlanInfo() types.PlanInfo {
	return &planInfo{s.summary.GetPlanInfo()}
}

func (s *summary) QueryStats() types.QueryStats {
	return &queryStats{s.summary.GetQueryStats()}
}

func (p *planInfo) Id() string {
	return string(p.planInfo.Id)
}

func (p *planInfo) Name() string {
	return string(p.planInfo.Name)
}

func (p *planInfo) Details() string {
	return string(p.planInfo.Details)
}

func (p *planInfo) Columns() []string {
	columnsVec := make([]string, 0, len(p.planInfo.Columns))
	for _, column := range p.planInfo.Columns {
		columnsVec = append(columnsVec, string(column))
	}
	return columnsVec
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

func (p *planInfo) Children() []types.PlanInfo {
	children := make([]types.PlanInfo, 0, len(p.planInfo.Children))
	for _, child := range p.planInfo.Children {
		children = append(children, &planInfo{planInfo: child})
	}
	return children
}

func (q *queryStats) NumAffectedNodes() int64 {
	return q.queryStats.NumAffectedNodes
}

func (q *queryStats) NumAffectedEdges() int64 {
	return q.queryStats.NumAffectedEdges
}
