package main

import (
	"encoding/json"
	"fmt"
)

func (s EvalSpec) Expr() Expr {
	if s.Commit != nil {
		return *s.Commit
	}
	if s.Tree != nil {
		return *s.Tree
	}
	if s.Blob != nil {
		return *s.Blob
	}
	panic(fmt.Errorf("unset spec"))
}

// Spec can be recursive
// Val isn't
func Parse(specText string) (EvalSpec, error) {
	var spec EvalSpec
	err := json.Unmarshal([]byte(specText), &spec)
	if (spec.Commit == nil && spec.Tree == nil && spec.Blob == nil) || spec.SHA == "" {
		return spec, fmt.Errorf("could not parse spec from %q", specText)
	}
	return spec, err
}
