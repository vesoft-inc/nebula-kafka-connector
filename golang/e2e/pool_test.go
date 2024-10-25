package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

func TestPoolSessionSet(t *testing.T) {
	params := map[string]string{
		"s": "\"1\"",
		"i": "2",
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
	resp, err := c.Execute(`return zoned_datetime("2020-03-02T23:12:00+0000")`)
	if err != nil {
		t.Fatal(err)
	}
	var dt nebula.NullZonedDatetime
	if err := resp.Scan(&dt); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dt.Valid, true)
	assert.Equal(t, dt.Data.GetOffset(), 8*3600)
	assert.Equal(t, dt.Data.GetDay(), 3)
	assert.Equal(t, dt.Data.GetHour(), 7)
	assert.Equal(t, dt.Data.GetMinute(), 12)

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

func TestSessionSet(t *testing.T) {
	p, err := nebula.NewNebulaPool(nebulaAddress, nebulaUser, nebulaPassword,
		nebula.WithPoolGraph("test_graph"),
		nebula.WithPoolSchema("/test_schema"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	c, err := p.GetClient()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Execute(`show current_schema`)
	if err != nil {
		t.Fatal(err)
	}
	var name, path, owner nebula.NullString
	if err := resp.Scan(&name, &path, &owner); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, path.Valid, true)
	assert.Equal(t, string(path.Data), "/test_schema")
	assert.Equal(t, owner.Valid, true)
	assert.Equal(t, string(owner.Data), "root")
	resp, err = c.Execute(`show current_session`)
	if err != nil {
		t.Fatal(err)
	}
	row, err := resp.Next()
	if err != nil {
		t.Fatal(err)
	}
	v, err := row.GetValueByName("home_graph_name")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, v.String(), "test_graph")
}
