package conf

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {

	mockConf := `
[ds]
	[ds.redis]
		port = "6379"
		chan = "orders"
[srv]
	port = "1234"
`

	r := strings.NewReader(mockConf)
	app, err := Read(r)
	if err != nil {
		t.Errorf("Read error: %s", err)
	}
	want := "6379"
	if got := app.DS.Redis.Port; got != want {
		t.Errorf("Datastore Redis Port %d, want %d", got, want)
	}

	want = "1234"
	t.Logf("%+v\n %+v\n", app.SRV.Port, want)
	/*
		if got := app.SRV.Port; got != want {
			t.Log(got, want)
			t.Errorf("Websocket Port %d, want %d", got, want)
		}
	*/
}
