package core

type Model interface {
	Save() (bool, error)
}
