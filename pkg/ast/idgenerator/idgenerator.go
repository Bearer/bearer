package idgenerator

type Generator struct {
	next uint32
}

func New() *Generator {
	return &Generator{}
}

func (generator *Generator) Get() uint32 {
	next := generator.next
	generator.next++
	return next
}
