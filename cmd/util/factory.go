package util

import (
	"github.com/ForbiddenR/auto-commit/pkg/resource"
)

type Factory interface {
	NewBuilder() *resource.Builder
}

type factoryImpl struct{}

func NewFactory() Factory {
	return &factoryImpl{}
}

func (f *factoryImpl) NewBuilder() *resource.Builder {
	return resource.NewBuilder()
}
