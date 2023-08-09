package resource

type Visitor interface {
	Visit() error
}

type VisitorFunc func(error) error
