// Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"bytes"
	"os"
	"testing"
)

func TestApplyHerokuSpecificOverrides(t *testing.T) {
	fixtures := map[string]string{
		"DATABASE_URL": "https://my.postgres:5432",
		"PORT":         "9999",
		"APP_SECRET":   "abcdefghijk",
	}
	for key, value := range fixtures {
		defer os.Setenv(key, os.Getenv(key))
		os.Setenv(key, value)
	}

	c := &Config{}
	c.Server.Port = 5555
	c.Database.ConnectionString = "/tmp/db.sqlite"
	c.Secret = []byte("123456789")
	c.Server.AutoTLS = []string{"www.khulnasoft.com"}

	if err := applyHerokuSpecificOverrides(c); err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if c.Database.ConnectionString != "https://my.postgres:5432" {
		t.Errorf("Unexpected connection string value %v", c.Database.ConnectionString)
	}

	if c.Server.Port != 9999 {
		t.Errorf("Unexpected port value %v", c.Server.Port)
	}

	if !bytes.Equal(c.Secret.Bytes(), []byte("abcdefghijk")) {
		t.Errorf("Unexpected cookie secret %v", c.Secret.Bytes())
	}

	if c.Server.AutoTLS[0] != "www.khulnasoft.com" {
		t.Errorf("Unexpected AutoTLS config %v", c.Server.AutoTLS)
	}
}
