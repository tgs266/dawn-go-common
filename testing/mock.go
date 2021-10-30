package testing

import (
	"log"

	"github.com/undefinedlabs/go-mpatch"
)

type Mock struct {
	Patches []*mpatch.Patch
}

func CreateMock(target, redirection interface{}) *Mock {
	patch, err := mpatch.PatchMethod(target, redirection)
	if err != nil {
		log.Fatal(err)
	}

	var patches []*mpatch.Patch
	patches = append(patches, patch)

	mock := &Mock{
		Patches: patches,
	}
	return mock
}

func (m *Mock) AddMock(target, redirection interface{}) *Mock {
	patch, err := mpatch.PatchMethod(target, redirection)
	if err != nil {
		log.Fatal(err)
	}

	m.Patches = append(m.Patches, patch)
	return m
}

func (m *Mock) Unpatch() {
	for i := 0; i < len(m.Patches); i++ {
		m.Patches[i].Unpatch()
	}
}
