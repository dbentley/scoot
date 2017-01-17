package main

type CounterModel struct {
	Ops []Update
}

type Update struct {
	Old string
	New string
}
