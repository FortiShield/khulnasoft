// Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
// SPDX-License-Identifier: Apache-2.0

package keys

import "errors"

// different errors will be returned for different validation failures
var (
	ErrPasswordTooShort = errors.New("keys: given password is shorter than 8 characters")
	ErrPasswordTooLong  = errors.New("keys: given password is longer than 64 characters")
)

// ValidatePassword checks whether the given password meets all requirements
// of the currently applicable password policy
func ValidatePassword(pw string) error {
	if len(pw) < 8 {
		return ErrPasswordTooShort
	}
	if len(pw) > 64 {
		return ErrPasswordTooLong
	}
	return nil
}
