// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package yttlibrary

import (
	"fmt"

	"github.com/k14s/starlark-go/starlark"
	"github.com/k14s/starlark-go/starlarkstruct"
	"github.com/vmware-tanzu/carvel-ytt/pkg/template/core"
)

var (
	AssertAPI = starlark.StringDict{
		"assert": &starlarkstruct.Module{
			Name: "assert",
			Members: starlark.StringDict{
				"equals": starlark.NewBuiltin("assert.equals", core.ErrWrapper(assertModule{}.Equals)),
				"fail":   starlark.NewBuiltin("assert.fail", core.ErrWrapper(assertModule{}.Fail)),
				"try_to": starlark.NewBuiltin("assert.try_to", core.ErrWrapper(assertModule{}.TryTo)),
				//"min_length": /*extract in var */ starlark.NewBuiltin("assert.min_len", core.ErrWrapper(assertModule{}.MinLength)),
			},
		},
	}
)

type assertModule struct{}

// Equals compares two values for equality
func (b assertModule) Equals(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if args.Len() != 2 {
		return starlark.None, fmt.Errorf("expected two arguments")
	}

	expected := args.Index(0)
	if _, notOk := expected.(starlark.Callable); notOk {
		return starlark.None, fmt.Errorf("expected argument not to be a function, but was %T", expected)
	}

	actual := args.Index(1)
	if _, notOk := actual.(starlark.Callable); notOk {
		return starlark.None, fmt.Errorf("expected argument not to be a function, but was %T", actual)
	}

	expectedString, err := b.asString(expected)
	if err != nil {
		return starlark.None, err
	}

	actualString, err := b.asString(actual)
	if err != nil {
		return starlark.None, err
	}

	if expectedString != actualString {
		return starlark.None, fmt.Errorf("Not equal:\n"+
			"(expected type: %s)\n%s\n\n(was type: %s)\n%s", expected.Type(), expectedString, actual.Type(), actualString)
	}

	return starlark.None, nil
}

func (b assertModule) asString(value starlark.Value) (string, error) {
	starlarkValue, err := core.NewStarlarkValue(value).AsGoValue()
	if err != nil {
		return "", err
	}
	yamlString, err := yamlModule{}.Encode(starlarkValue)
	if err != nil {
		return "", err
	}
	return yamlString, nil
}

func (b assertModule) Fail(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if args.Len() != 1 {
		return starlark.None, fmt.Errorf("expected exactly one argument")
	}

	val, err := core.NewStarlarkValue(args.Index(0)).AsString()
	if err != nil {
		return starlark.None, err
	}

	return starlark.None, fmt.Errorf("fail: %s", val)
}

func (b assertModule) TryTo(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if args.Len() != 1 {
		return starlark.None, fmt.Errorf("expected exactly one argument")
	}

	lambda := args.Index(0)
	if _, ok := lambda.(starlark.Callable); !ok {
		return starlark.None, fmt.Errorf("expected argument to be a function, but was %T", lambda)
	}

	retVal, err := starlark.Call(thread, lambda, nil, nil)
	if err != nil {
		return starlark.Tuple{starlark.None, starlark.String(err.Error())}, nil
	}
	return starlark.Tuple{retVal, starlark.None}, nil
}

//func (b assertModule) MinLength(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
//	//do we need two values here?
//	if args.Len() != 1 {
//		return starlark.None, fmt.Errorf("expected exactly one argument")
//	}
//
//	//convert string function to
//	return starlark.None, nil
//}

func newMinLengthStarlarkFunc(minLength int) core.StarlarkFunc {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if args.Len() != 1 {
			return starlark.None, fmt.Errorf("expected exactly one argument")
		}
		valLen := starlark.Len(args[0])
		if valLen < 0 {
			return starlark.None, fmt.Errorf("expected something that had length")
		}
		if valLen >= minLength {
			return starlark.None, nil
		} else {
			return starlark.None, fmt.Errorf("length of value was less than %v", minLength)
		}
	}
}
func NewAssertMinLength(minLength int) starlark.Callable {
	return starlark.NewBuiltin("assert.min_len", core.ErrWrapper(newMinLengthStarlarkFunc(minLength)))
}

func newMaxLengthStarlarkFunc(maxLength int) core.StarlarkFunc {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if args.Len() != 1 {
			return starlark.None, fmt.Errorf("expected exactly one argument")
		}
		valLen := starlark.Len(args[0])
		if valLen < 0 {
			return starlark.None, fmt.Errorf("expected something that had length")
		}
		if valLen <= maxLength {
			return starlark.None, nil
		} else {
			return starlark.None, fmt.Errorf("length of value was less than %v", maxLength)
		}
	}
}
func NewAssertMaxLength(maxLength int) starlark.Callable {
	return starlark.NewBuiltin("assert.max_len", core.ErrWrapper(newMaxLengthStarlarkFunc(maxLength)))
}

func newMinStarlarkFunc(min int) core.StarlarkFunc {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if args.Len() != 1 {
			return starlark.None, fmt.Errorf("expected exactly one argument")
		}
		val := args[0]
		v, err := starlark.NumberToInt(val)
		if err != nil {
			return starlark.None, fmt.Errorf("expected value to be an number, but was %s", val.Type())
		}
		num, _ := v.Int64()
		intNum := int(num)
		if intNum >= min {
			return starlark.None, nil
		} else {
			return starlark.None, fmt.Errorf("value was less than %v", min)
		}
	}
}
func NewAssertMin(min int) starlark.Callable {
	return starlark.NewBuiltin("assert.min", core.ErrWrapper(newMinStarlarkFunc(min)))
}

func newMaxStarlarkFunc(max int) core.StarlarkFunc {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if args.Len() != 1 {
			return starlark.None, fmt.Errorf("expected exactly one argument")
		}
		val := args[0]
		v, err := starlark.NumberToInt(val)
		if err != nil {
			return starlark.None, fmt.Errorf("expected value to be an number, but was %s", val.Type())
		}
		num, _ := v.Int64()
		intNum := int(num)
		if intNum <= max {
			return starlark.None, nil
		} else {
			return starlark.None, fmt.Errorf("length of value was less than %v", max)
		}
	}
}
func NewAssertMax(max int) starlark.Callable {
	return starlark.NewBuiltin("assert.max", core.ErrWrapper(newMaxStarlarkFunc(max)))
}
