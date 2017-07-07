package main

import (
	"io/ioutil"
	"testing"

	"github.com/convox/praxis/sdk/rack"
	"github.com/convox/praxis/types"
	"github.com/stretchr/testify/assert"
)

func TestAppCreate(t *testing.T) {
	Rack, err := rack.NewFromEnv()
	assert.NoError(t, err)

	app, err := Rack.AppCreate("")
	assert.EqualError(t, err, "app name required")
	assert.Nil(t, app)

	app, err = Rack.AppCreate("3")
	assert.EqualError(t, err, "app name invalid")
	assert.Nil(t, app)

	app, err = appCreate(Rack, "valid")
	defer appDelete(Rack, app.Name)
	assert.NoError(t, err)
	assert.EqualValues(t, &types.App{
		Name:    "valid",
		Release: "",
		Status:  "running",
	}, app)

	app, err = Rack.AppCreate("valid")
	assert.EqualError(t, err, "app already exists: valid")
	assert.Nil(t, app)
}

func TestAppDelete(t *testing.T) {
	Rack, err := rack.NewFromEnv()
	assert.NoError(t, err)

	err = Rack.AppDelete("")
	assert.EqualError(t, err, "response status 404")

	err = Rack.AppDelete("3")
	assert.EqualError(t, err, "app name invalid")

	err = Rack.AppDelete("missing")
	assert.EqualError(t, err, "no such app: missing")

	_, err = Rack.AppCreate("valid")
	assert.NoError(t, err)
	err = Rack.AppDelete("valid")
	assert.NoError(t, err)
}

func TestAppGet(t *testing.T) {
	Rack, err := rack.NewFromEnv()
	assert.NoError(t, err)

	app, err := Rack.AppGet("")
	assert.EqualError(t, err, "response status 404")
	assert.Nil(t, app)

	app, err = Rack.AppGet("3")
	assert.EqualError(t, err, "app name invalid")
	assert.Nil(t, app)

	app, err = appCreate(Rack, "valid")
	defer appDelete(Rack, app.Name)
	assert.NoError(t, err)
	a, err := Rack.AppGet("valid")
	assert.NoError(t, err)
	assert.EqualValues(t, app, a)
}

func TestAppList(t *testing.T) {
	Rack, err := rack.NewFromEnv()
	assert.NoError(t, err)

	apps, err := Rack.AppList()
	assert.NoError(t, err)
	assert.EqualValues(t, types.Apps{}, apps)

	app, err := appCreate(Rack, "foo")
	defer appDelete(Rack, app.Name)
	assert.NoError(t, err)

	app, err = appCreate(Rack, "bar")
	defer appDelete(Rack, app.Name)
	assert.NoError(t, err)

	apps, err = Rack.AppList()
	assert.NoError(t, err)
	assert.EqualValues(t, types.Apps{
		types.App{
			Name:    "bar",
			Release: "",
			Status:  "running",
		},
		types.App{
			Name:    "foo",
			Release: "",
			Status:  "running",
		},
	}, apps)
}

func TestAppLogs(t *testing.T) {
	Rack, err := rack.NewFromEnv()
	assert.NoError(t, err)

	app, err := appCreate(Rack, "valid")
	defer appDelete(Rack, app.Name)

	r, err := Rack.AppLogs(app.Name, types.LogsOptions{})
	assert.NoError(t, err)
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, []byte{}, b)

	// FIXME: assert app process logs
	// FIXME: assert log options filter, follow, prefix, since
}

func TestAppRegistry(t *testing.T) {
	Rack, err := rack.NewFromEnv()
	assert.NoError(t, err)

	app, err := appCreate(Rack, "valid")
	defer appDelete(Rack, app.Name)

	r, err := Rack.AppRegistry(app.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, r.Hostname)
}
