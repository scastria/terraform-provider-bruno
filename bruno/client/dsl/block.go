package dsl

type BruBlock interface {
	GetTag() string
	Export() string
}
