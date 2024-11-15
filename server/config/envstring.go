// Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
// SPDX-License-Identifier: Apache-2.0

package config

// EnvString is a string that expands environemt variables
type EnvString string

// Decode validates and assigns v.
func (e *EnvString) Decode(v string) error {
	*e = EnvString(v)
	return nil
}

func (e *EnvString) String() string {
	return ExpandString(string(*e))
}

// RawString returns the original value without any interpolation.
func (e *EnvString) RawString() string {
	return string(*e)
}
