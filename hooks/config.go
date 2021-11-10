package hooks

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Okta struct {
	CLIENT struct {
		ORG_URL string
		TOKEN   string
	}
	HOOK struct {
		BASE_URL string
	}
}

type Configuration struct {
	Okta
}

var Config Configuration

func init() {
	buf, err := ioutil.ReadFile(".okta.yaml")
	if err != nil {
		log.Fatalf("No Configuration file exists yet: %v\n", err)
	}

	err = yaml.Unmarshal(buf, &Config)
	if err != nil {
		log.Fatal(err)
	}

}
