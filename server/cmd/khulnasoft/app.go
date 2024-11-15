// Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/khulnasoft/khulnasoft/server/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type app struct {
	logger *logrus.Logger
	config *config.Config
}

func newApp(populateMissing, quiet bool, envFileOverride string) *app {
	logger := logrus.New()
	cfg, cfgErr := config.New(populateMissing, envFileOverride)
	if cfgErr != nil {
		if errors.Is(cfgErr, config.ErrPopulatedMissing) {
			logger.Infof("Some configuration values were missing: %v", cfgErr.Error())
		} else {
			logger.WithError(cfgErr).Fatal("Error sourcing runtime configuration")
		}
	}

	logger.SetLevel(cfg.App.LogLevel.LogLevel())
	if !quiet && !cfg.SMTPConfigured() {
		logger.Warn("SMTP for transactional email is not configured right now, mail delivery will be unreliable")
		logger.Warn("Refer to the documentation to find out how to configure SMTP")
	}
	return &app{
		logger: logger,
		config: cfg,
	}
}

func newLogger() *logrus.Logger {
	return logrus.New()
}

func newDB(c *config.Config, l *logrus.Logger) (*gorm.DB, error) {
	var d gorm.Dialector
	switch c.Database.Dialect.String() {
	case "sqlite3":
		d = sqlite.Open(c.Database.ConnectionString.String())
	case "mysql":
		d = mysql.Open(c.Database.ConnectionString.String())
	case "postgres":
		d = postgres.Open(c.Database.ConnectionString.String())
	}

	logLevel := logger.Silent
	if c.App.Development {
		logLevel = logger.Info
	}

	var gormDB *gorm.DB
	if err := backoff.RetryNotify(
		func() error {
			var err error
			gormDB, err = gorm.Open(d, &gorm.Config{
				Logger:                                   logger.Default.LogMode(logLevel),
				DisableForeignKeyConstraintWhenMigrating: c.Database.Dialect.String() == "sqlite3",
			})
			return err
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(c.Database.ConnectionRetries)),
		func(err error, duration time.Duration) {
			if l != nil && c.Database.ConnectionRetries != 0 {
				l.WithError(err).Warn("Connecting to database failed")
				l.WithField("duration", duration).Info("Scheduling sleep before retrying")
			}
		},
	); err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if c.Database.Dialect == "sqlite3" {
		db, err := gormDB.DB()
		if err != nil {
			return nil, fmt.Errorf("error accessing underlying database: %w", err)
		}
		db.SetMaxOpenConns(1)
	}
	return gormDB, nil
}
