package e2e

import (
	"testing"

	"github.com/go-playground/assert/v2"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

func TestPoolSessionSet(t *testing.T) {
	params := map[string]string{
		"$s": "\"1\"",
		"$i": "2",
	}
	p, err := nebula.NewNebulaPool(nebulaAddress, nebulaUser, nebulaPassword,
		nebula.WithPoolTimezone("Asia/Shanghai"),
		nebula.WithPoolParameters(params),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	c, err := p.GetClient()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Execute(`return zoned_datetime("2020-03-02T01:00:00+0000")`)
	if err != nil {
		t.Fatal(err)
	}
	var dt nebula.NullZonedDatetime
	if err := resp.Scan(&dt); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dt.Valid, true)
	assert.Equal(t, dt.Data.GetOffset(), 8*3600)

	resp, err = c.Execute(`return $s, $i`)
	if err != nil {
		t.Fatal(err)
	}
	var s nebula.NullString
	var i nebula.NullInt
	if err := resp.Scan(&s, &i); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, s.Valid, true)
	assert.Equal(t, string(s.Data), "1")
	assert.Equal(t, i.Valid, true)
	assert.Equal(t, int(i.Data), 2)
}
