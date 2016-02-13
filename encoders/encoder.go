package encoders

// an interface that defines a start function, used for setup.
type BaseEncoder interface {
	Start() error
}
