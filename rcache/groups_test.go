package rcache

import "testing"

func TestPutGet(t *testing.T) {
	Grouplist.Put("test", Group{Name: "TestName", FleetID: 202})
	_, err := Grouplist.Get("test")
	if err != nil {
		t.Error(err)
	}
}
