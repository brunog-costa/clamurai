package hashing

type Hasher struct {
	algo string
}

type HasherWorker interface {
}

func New() *Hasher                         {}
func (h *Hasher) Hashsum() (string, error) {}
