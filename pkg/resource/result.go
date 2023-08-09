package resource

type Result struct {
	err     error
	visitor Visitor
}

// Err returns the first error that occurred during the visit.
func (r *Result) Err() error {
	return r.err
}

func (r *Result) Visitor() Visitor {
	return r.visitor
}
