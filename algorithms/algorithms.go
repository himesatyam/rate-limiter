package algorithms

type Algorithms interface {
	Allow(key string) (bool, error)
}
