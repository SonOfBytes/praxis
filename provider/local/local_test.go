package local_test

import (
	"io/ioutil"
	"os"

	"github.com/convox/praxis/provider/local"
)

func testProvider() (*local.Provider, error) {
	tmp, err := ioutil.TempDir("", "praxis")
	if err != nil {
		return nil, err
	}

	return &local.Provider{Root: tmp, Test: true}, nil
}

func cleanup(p *local.Provider) {
	if p.Root != "" {
		os.RemoveAll(p.Root)
	}
}
