package config

type Google struct {
	ServiceAccount string `required:"true" split_words:"true"`
}
