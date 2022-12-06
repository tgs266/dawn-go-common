package optional

type Optional[T any] struct {
	err  error
	data T
}

func (o *Optional[T]) Get() T {
	return o.data
}

func (o *Optional[T]) GetOrPanic() T {
	if o.err != nil {
		panic(o.err)
	}
	return o.data
}

func (o *Optional[T]) GetError() error {
	return o.err
}
