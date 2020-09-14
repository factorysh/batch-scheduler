package action

import (
	"fmt"

	cs "github.com/compose-spec/compose-go/loader"
	"gopkg.in/yaml.v2"
)

// Compose represent a struct containing a raw docker-compose.yml file
type Compose struct {
	Raw    string
	Parsed map[string]interface{}
}

// NewCompose inits a new compose file struct
func NewCompose(input []byte) Compose {
	return Compose{
		Raw: string(input),
	}

}

// Parse ensures a docker-compose file, ensure content is valid
func (c *Compose) Parse() error {

	parsed, err := cs.ParseYAML([]byte(c.Raw))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error when validating compose file: %v", err))
	}

	c.Parsed = parsed

	return nil
}

// Recompose rewrite the file back
func (c *Compose) Recompose() (string, error) {

	ret, err := yaml.Marshal(&c.Parsed)
	if err != nil {
		return "", err
	}

	return string(ret), nil

}