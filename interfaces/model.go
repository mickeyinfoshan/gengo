package interfaces

type Model interface {
	Save() error
	Delete() error
}
