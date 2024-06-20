package yamlembed

import (
	"strings"
)

type Foo struct {
	A string `yaml:"aa"`
	p int64  `yaml:"-"`
}

type Bar struct {
	I      int64    `yaml:"-"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"`
	OI     []string `yaml:"oi,omitempty"`
	F      []any    `yaml:"f,flow"`
}

func (b *Bar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Bar
	if err := unmarshal((*plain)(b)); err != nil {
		return err
	}

	b.UpperB = strings.ToUpper(b.B)
	return nil
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}

func (b *Baz) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var foo Foo
	var bar Bar

	if err := unmarshal(&foo); err != nil {
		return err
	}

	if err := unmarshal(&bar); err != nil {
		return err
	}

	b.Foo = foo
	b.Bar = bar
	return nil
}
