package rsm

type Transition interface {
	Transition()
}

type Result interface {
	Result()
}

type Machine interface {
	Empty() Model
	Encode(Transition) (string, error)
	Decode(string) (Transition, error)
}

type Model interface {
	Copy() Model
	Apply(t Transition) Result
}
