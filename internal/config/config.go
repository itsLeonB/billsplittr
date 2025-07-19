package config

type Google struct {
	ServiceAccount string `split_words:"true" required:"true"`
}

func (g Google) Prefix() string {
	return "GOOGLE"
}
