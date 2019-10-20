package infra

type Specification struct {
	Debug              bool
	Port               int      `default:"4000"`
	SupportedLanguages []string `envconfig:"SUPPORTED_LANNGUAGES"`
}
