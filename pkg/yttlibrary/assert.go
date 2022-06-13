// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package yttlibrary

import (
	"fmt"

	"github.com/k14s/starlark-go/syntax"

	"github.com/k14s/starlark-go/starlark"
	"github.com/k14s/starlark-go/starlarkstruct"
	"github.com/vmware-tanzu/carvel-ytt/pkg/template/core"
)

var (
	AssertAPI = starlark.StringDict{
		"assert": &starlarkstruct.Module{
			Name: "assert",
			Members: starlark.StringDict{
				"equals":     starlark.NewBuiltin("assert.equals", core.ErrWrapper(assertModule{}.Equals)),
				"fail":       starlark.NewBuiltin("assert.fail", core.ErrWrapper(assertModule{}.Fail)),
				"try_to":     starlark.NewBuiltin("assert.try_to", core.ErrWrapper(assertModule{}.TryTo)),
				"min":        starlark.NewBuiltin("assert.min", core.ErrWrapper(assertModule{}.Min)),
				"min_length": starlark.NewBuiltin("assert.min_len", core.ErrWrapper(assertModule{}.MinLength)),
				"max":        starlark.NewBuiltin("assert.max", core.ErrWrapper(assertModule{}.Max)),
				"max_length": starlark.NewBuiltin("assert.max_length", core.ErrWrapper(assertModule{}.MaxLength)),
				"not_null":   starlark.NewBuiltin("assert.not_null", core.ErrWrapper(assertModule{}.NotNull)),
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
func (b assertModule) MinLength(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if args.Len() != 1 {
		return starlark.None, fmt.Errorf("expected exactly one argument")
	}

	// convert string function to
	return starlark.None, nil
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
			return starlark.None, fmt.Errorf("length of value was more than %v", maxLength)
		}
	}
}
func NewAssertMaxLength(maxLength int) starlark.Callable {
	return starlark.NewBuiltin("assert.max_len", core.ErrWrapper(newMaxLengthStarlarkFunc(maxLength)))
}
func (b assertModule) MaxLength(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, nil
}

// assertMinimum produces a higher-order Starlark function that asserts that a given value is at least "minimum".
func assertMinimum(minimum starlark.Value) (*starlark.Function, error) {
	src := `lambda value: fail("{} is less than {}".format(value, minimum)) if value < minimum else None`
	expr, err := syntax.ParseExpr("@ytt:assert.min()", src, syntax.BlockScanner)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse internal expression (%s) :%s", src, err))
	}
	thread := &starlark.Thread{Name: "ytt-internal"}

	evalExpr, err := starlark.EvalExpr(thread, expr, starlark.StringDict{"minimum": minimum})
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke @ytt:assert.min(%v) :%s", minimum, err)
	}
	return evalExpr.(*starlark.Function), nil
}

// NewAssertMin produces a higher-order Starlark function that asserts that a given value is at least "minimum"
func NewAssertMin(minimum int) *starlark.Function {
	min := core.NewGoValue(minimum).AsStarlarkValue()
	minimumFunc, err := assertMinimum(min)
	if err != nil {
		// TODO: given that "minimum" is user-supplied, return "err" instead of panicing
		panic(fmt.Sprintf("failed to build assert.minimum(): %s", err))
	}
	return minimumFunc
}

// Min is a core.StarlarkFunc that asserts that a given value is at least a given minimum.
func (b assertModule) Min(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if len(args) == 0 {
		return starlark.None, fmt.Errorf("expected at least one argument.")
	}
	if len(args) > 2 {
		return starlark.None, fmt.Errorf("expected at no more than two arguments.")
	}

	minFunc, err := assertMinimum(args[0])
	if err != nil {
		return starlark.None, err
	}
	if len(args) == 1 {
		return minFunc, nil
	}

	result, err := starlark.Call(thread, minFunc, starlark.Tuple{args[1]}, []starlark.Tuple{})
	if err != nil {
		return starlark.None, err
	}
	return result, nil
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
			return starlark.None, fmt.Errorf("value was more than %v", max)
		}
	}
}
func NewAssertMax(max int) starlark.Callable {
	return starlark.NewBuiltin("assert.max", core.ErrWrapper(newMaxStarlarkFunc(max)))
}
func (b assertModule) Max(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, nil
}

func newNotNullStarlarkFunc() core.StarlarkFunc {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if args.Len() != 1 {
			return starlark.None, fmt.Errorf("expected exactly one argument")
		}
		_, ok := args[0].(starlark.NoneType)
		if ok {
			return starlark.None, fmt.Errorf("value was null")
		}

		return starlark.None, nil
	}
}
func NewAssertNotNull() starlark.Callable {
	return starlark.NewBuiltin("assert.not_null", core.ErrWrapper(newNotNullStarlarkFunc()))
}
func (b assertModule) NotNull(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, nil
}
