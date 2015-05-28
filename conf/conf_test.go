package conf

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {

	mockConf := `
[ds]
    [ds.redis]
		host = ":6379"
[srv]
    port = "1234"
[log]
    path = "info.log"
[errors]
    [errors.NotExistInCache]
    msg = "such record does not exist"
`

	r := strings.NewReader(mockConf)
	app, err := Read(r)
	if err != nil {
		t.Errorf("Read error: %s", err)
	}

	want := ":6379"
	if got := app.DS.Redis.Host; got != want {
		t.Errorf("Datastore Redis Host %d, want %d", got, want)
	}

	want = "1234"
	if got := app.SRV.Port; got != want {
		t.Errorf("Server Port %d, want %d", got, want)
	}

	want = "info.log"
	if got := app.Log.Path; got != want {
		t.Errorf("got %s, want %s", got, want)
	}

	want = "such record does not exist"
	if got := app.ErrorMsg["NotExistInCache"].Msg; got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
