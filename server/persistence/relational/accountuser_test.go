// Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
// SPDX-License-Identifier: Apache-2.0

package relational

import (
	"fmt"
	"reflect"
	"testing"

	"gorm.io/gorm"
	"github.com/khulnasoft/khulnasoft/server/persistence"
)

func TestRelationalDAL_CreateAccountUser(t *testing.T) {
	tests := []struct {
		name        string
		arg         *persistence.AccountUser
		expectError bool
		assertion   dbAccess
	}{
		{
			"ok",
			&persistence.AccountUser{
				AccountUserID: "account-user-id",
			},
			false,
			func(db *gorm.DB) error {
				if err := db.Where("account_user_id = ?", "account-user-id").First(&AccountUser{}).Error; err != nil {
					return fmt.Errorf("error looking up account user: %w", err)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, closeDB := createTestDatabase()
			defer closeDB()

			dal := NewRelationalDAL(db)

			err := dal.CreateAccountUser(test.arg)
			if (err != nil) != test.expectError {
				t.Errorf("Unexpected error value %v", err)
			}

			if err := test.assertion(db); err != nil {
				t.Errorf("Assertion error validating database content: %v", err)
			}
		})
	}
}

func TestRelationalDAL_FindAccountUser(t *testing.T) {
	tests := []struct {
		name           string
		setup          dbAccess
		query          interface{}
		expectedResult persistence.AccountUser
		expectError    bool
	}{
		{
			"bad query",
			noop,
			complex128(12),
			persistence.AccountUser{},
			true,
		},
		{
			"by user id found - include relationships",
			func(db *gorm.DB) error {
				if err := db.Save(&AccountUser{
					AccountUserID: "user-id",
					HashedEmail:   "xyz123",
				}).Error; err != nil {
					return fmt.Errorf("error saving fixture data: %v", err)
				}
				if err := db.Save(&AccountUserRelationship{
					AccountUserID:                     "user-id",
					AccountID:                         "account-id",
					RelationshipID:                    "relationship-id",
					PasswordEncryptedKeyEncryptionKey: "key",
				}).Error; err != nil {
					return fmt.Errorf("error saving fixture data: %v", err)
				}
				return nil
			},
			persistence.FindAccountUserQueryByAccountUserIDIncludeRelationships("user-id"),
			persistence.AccountUser{
				AccountUserID: "user-id",
				HashedEmail:   "xyz123",
				Relationships: []persistence.AccountUserRelationship{
					{
						AccountUserID:                     "user-id",
						AccountID:                         "account-id",
						RelationshipID:                    "relationship-id",
						PasswordEncryptedKeyEncryptionKey: "key",
					},
				},
			},
			false,
		},
		{
			"by user id not found - include relationships",
			func(db *gorm.DB) error {
				if err := db.Save(&AccountUser{
					AccountUserID: "user-id",
					HashedEmail:   "xyz123",
				}).Error; err != nil {
					return fmt.Errorf("error saving fixture data: %v", err)
				}
				if err := db.Save(&AccountUserRelationship{
					AccountUserID:  "user-id",
					AccountID:      "account-id",
					RelationshipID: "relationship-id",
				}).Error; err != nil {
					return fmt.Errorf("error saving fixture data: %v", err)
				}
				return nil
			},
			persistence.FindAccountUserQueryByAccountUserIDIncludeRelationships("user-id-2"),
			persistence.AccountUser{},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, closeDB := createTestDatabase()
			defer closeDB()

			dal := NewRelationalDAL(db)

			if err := test.setup(db); err != nil {
				t.Fatalf("Error setting up test: %v", err)
			}

			result, err := dal.FindAccountUser(test.query)

			if !reflect.DeepEqual(test.expectedResult, result) {
				t.Errorf("Expected %v, got %v", test.expectedResult, result)
			}

			if (err != nil) != test.expectError {
				t.Errorf("Unexpected error value %v", err)
			}
		})
	}
}

func TestRelationalDAL_UpdateAccountUser(t *testing.T) {
	tests := []struct {
		name        string
		setup       dbAccess
		arg         *persistence.AccountUser
		expectError bool
		assertion   dbAccess
	}{
		{
			"user does not exist",
			func(db *gorm.DB) error {
				if err := db.Save(&AccountUser{
					AccountUserID:  "account-user-id-z",
					HashedPassword: "abc123",
				}).Error; err != nil {
					return fmt.Errorf("error creating account user fixture: %v", err)
				}
				return nil
			},
			&persistence.AccountUser{
				AccountUserID:  "account-user-id-a",
				HashedPassword: "xyz987",
			},
			true,
			func(db *gorm.DB) error {
				var accountUser AccountUser
				if err := db.Where("account_user_id = ?", "account-user-id-z").First(&accountUser).Error; err != nil {
					return fmt.Errorf("error looking up record: %v", err)
				}
				if accountUser.HashedPassword != "abc123" {
					return fmt.Errorf("record unexpectedly changed with password %v", accountUser.HashedPassword)
				}
				return nil
			},
		},
		{
			"ok",
			func(db *gorm.DB) error {
				if err := db.Save(&AccountUser{
					AccountUserID:  "account-user-id-z",
					HashedPassword: "abc123",
				}).Error; err != nil {
					return fmt.Errorf("error creating account user fixture: %v", err)
				}
				return nil
			},
			&persistence.AccountUser{
				AccountUserID:  "account-user-id-z",
				HashedPassword: "xyz987",
			},
			false,
			func(db *gorm.DB) error {
				var accountUser AccountUser
				if err := db.Where("account_user_id = ?", "account-user-id-z").First(&accountUser).Error; err != nil {
					return fmt.Errorf("error looking up record: %v", err)
				}
				if accountUser.HashedPassword != "xyz987" {
					return fmt.Errorf("record not updated with password %v", accountUser.HashedPassword)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, closeDB := createTestDatabase()
			defer closeDB()

			dal := NewRelationalDAL(db)

			if err := test.setup(db); err != nil {
				t.Fatalf("Error setting up test: %v", err)
			}

			err := dal.UpdateAccountUser(test.arg)
			if (err != nil) != test.expectError {
				t.Errorf("Unexpected error value %v", err)
			}

			if err := test.assertion(db); err != nil {
				t.Errorf("Assertion error validating database content: %v", err)
			}
		})
	}
}

func TestRelationalDAL_FindAccountUsers(t *testing.T) {
	tests := []struct {
		name           string
		setup          dbAccess
		arg            interface{}
		expectError    bool
		expectedResult []persistence.AccountUser
	}{
		{
			"bad query",
			noop,
			"puppies",
			true,
			nil,
		},
		{
			"empty db",
			noop,
			persistence.FindAccountUsersQueryAllAccountUsers{},
			false,
			nil,
		},
		{
			"find users",
			func(db *gorm.DB) error {
				if err := db.Create(&AccountUser{
					AccountUserID: "account-user-a",
				}).Error; err != nil {
					return fmt.Errorf("error inserting fixture: %w", err)
				}
				if err := db.Create(&Account{
					Name:      "test",
					AccountID: "account-a",
				}).Error; err != nil {
					return fmt.Errorf("error inserting fixture: %w", err)
				}
				if err := db.Create(&AccountUserRelationship{
					RelationshipID:                    "relationship-a",
					AccountUserID:                     "account-user-a",
					AccountID:                         "account-a",
					PasswordEncryptedKeyEncryptionKey: "something",
				}).Error; err != nil {
					return fmt.Errorf("error inserting fixture: %w", err)
				}
				return nil
			},
			persistence.FindAccountUsersQueryAllAccountUsers{IncludeRelationships: true},
			false,
			[]persistence.AccountUser{
				{AccountUserID: "account-user-a", Relationships: []persistence.AccountUserRelationship{
					{RelationshipID: "relationship-a", AccountUserID: "account-user-a", AccountID: "account-a", PasswordEncryptedKeyEncryptionKey: "something"},
				}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, closeDB := createTestDatabase()
			defer closeDB()
			if err := test.setup(db); err != nil {
				t.Fatalf("Error setting up database %v", err)
			}

			dal := NewRelationalDAL(db)

			result, err := dal.FindAccountUsers(test.arg)
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error value %v", err)
			}

			if !reflect.DeepEqual(test.expectedResult, result) {
				t.Errorf("Expected %v, got %v", test.expectedResult, result)
			}
		})
	}
}
