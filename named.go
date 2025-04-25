// Copyright (c) 2025 Renorm Labs. All rights reserved.

package moduli

import (
	"reflect"
	"sync"
)

// Named associates a human-readable name with an option, used for mutation
// tracking. Has no effect on option behavior.
func Named[T any](name string, opt Option[T]) Option[T] {
	optionNames.Store(reflect.ValueOf(opt).Pointer(), name)
	return opt
}

var optionNames sync.Map

// optionName returns the registered name for an option, if available. Falls
// back to "option" if unnamed.
func optionName[T any](opt Option[T]) string {
	if v, ok := optionNames.Load(reflect.ValueOf(opt).Pointer()); ok {
		return v.(string)
	}
	return "option"
}
