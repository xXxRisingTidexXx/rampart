package logging

type Publication interface {
	URL() string
	Body() string
}
