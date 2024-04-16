package proto

import (
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/common"
)

func ConvertZonedTime(zdt *common.ZonedDatetime) *time.Time {
	if zdt == nil {
		return nil
	}
	timezone := time.FixedZone("", int(zdt.GetOffset()))
	t := time.Date(
		int(zdt.GetYear()),
		time.Month(zdt.GetMonth()),
		int(zdt.GetDay()),
		int(zdt.GetHour()),
		int(zdt.GetMinute()),
		int(zdt.GetSec()),
		int(zdt.GetMicrosec())*int(time.Microsecond),
		timezone,
	)
	return &t
}
