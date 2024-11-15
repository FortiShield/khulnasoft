// Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
// SPDX-License-Identifier: Apache-2.0

package persistence

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type assertion func(interface{}) error

type mockInsertEventDatabase struct {
	DataAccessLayer
	findAccountResult Account
	findAccountErr    error
	findSecretResult  Secret
	findSecretErr     error
	createEventErr    error
	methodArgs        []interface{}
}

func (m *mockInsertEventDatabase) FindAccount(q interface{}) (Account, error) {
	m.methodArgs = append(m.methodArgs, q)
	return m.findAccountResult, m.findAccountErr
}

func (m *mockInsertEventDatabase) FindSecret(q interface{}) (Secret, error) {
	m.methodArgs = append(m.methodArgs, q)
	return m.findSecretResult, m.findSecretErr
}

func (m *mockInsertEventDatabase) CreateEvent(e *Event) error {
	m.methodArgs = append(m.methodArgs, e)
	return m.createEventErr
}

func TestPersistenceLayer_Insert(t *testing.T) {
	tests := []struct {
		name           string
		callArgs       []string
		db             *mockInsertEventDatabase
		expectError    bool
		argsAssertions []assertion
	}{
		{
			"account lookup error",
			[]string{"user-id", "account-id", "payload"},
			&mockInsertEventDatabase{
				findAccountErr: errors.New("did not work"),
			},
			true,
			[]assertion{
				func(accountID interface{}) error {
					if cast, ok := accountID.(FindAccountQueryActiveByID); ok {
						if cast != "account-id" {
							return fmt.Errorf("unexpected account identifier %v", cast)
						}
					}
					return nil
				},
			},
		},
		{
			"user lookup error",
			[]string{"user-id", "account-id", "payload"},
			&mockInsertEventDatabase{
				findAccountResult: Account{
					Name:     "test",
					UserSalt: "{1,} CaHVhk78uhoPmf5wanA0vg==",
				},
				findSecretErr: errors.New("did not work"),
			},
			true,
			[]assertion{
				func(accountID interface{}) error {
					if cast, ok := accountID.(FindAccountQueryActiveByID); ok {
						if cast != "account-id" {
							return fmt.Errorf("unexpected account identifier %v", cast)
						}
					}
					return nil
				},
				func(userID interface{}) error {
					if cast, ok := userID.(FindSecretQueryBySecretID); ok {
						if cast == "user-id" || cast == "" {
							return fmt.Errorf("unexpected user identifier %v", cast)
						}
					}
					return nil
				},
			},
		},
		{
			"insert error",
			[]string{"user-id", "account-id", "payload"},
			&mockInsertEventDatabase{
				findAccountResult: Account{
					Name:     "test",
					UserSalt: "{1,} CaHVhk78uhoPmf5wanA0vg==",
				},
				createEventErr: errors.New("did not work"),
			},
			true,
			[]assertion{
				func(accountID interface{}) error {
					if cast, ok := accountID.(FindAccountQueryActiveByID); ok {
						if cast != "account-id" {
							return fmt.Errorf("unexpected account identifier %v", cast)
						}
					}
					return nil
				},
				func(userID interface{}) error {
					if cast, ok := userID.(FindSecretQueryBySecretID); ok {
						if cast == "user-id" || cast == "" {
							return fmt.Errorf("unexpected user identifier %v", cast)
						}
					}
					return nil
				},
				func(evt interface{}) error {
					if cast, ok := evt.(*Event); ok {
						wellformed := cast.Payload == "payload" &&
							cast.AccountID == "account-id" &&
							cast.EventID != "" &&
							*cast.SecretID != "user-id"
						if !wellformed {
							return fmt.Errorf("unexpected event shape %v", cast)
						}
					}
					return nil
				},
			},
		},
		{
			"ok",
			[]string{"user-id", "account-id", "payload"},
			&mockInsertEventDatabase{
				findAccountResult: Account{
					Name:     "test",
					UserSalt: "{1,} CaHVhk78uhoPmf5wanA0vg==",
				},
			},
			false,
			[]assertion{
				func(accountID interface{}) error {
					if cast, ok := accountID.(FindAccountQueryActiveByID); ok {
						if cast != "account-id" {
							return fmt.Errorf("unexpected account identifier %v", cast)
						}
					}
					return nil
				},
				func(userID interface{}) error {
					if cast, ok := userID.(FindSecretQueryBySecretID); ok {
						if cast == "user-id" || cast == "" {
							return fmt.Errorf("unexpected user identifier %v", cast)
						}
					}
					return nil
				},
				func(evt interface{}) error {
					if cast, ok := evt.(*Event); ok {
						wellformed := cast.Payload == "payload" &&
							cast.AccountID == "account-id" &&
							cast.EventID != "" &&
							*cast.SecretID != "user-id"
						if !wellformed {
							return fmt.Errorf("unexpected event shape %v", cast)
						}
					}
					return nil
				},
			},
		},
		{
			"anonymous event ok",
			[]string{"", "account-id", "payload"},
			&mockInsertEventDatabase{
				findAccountResult: Account{
					Name:     "test",
					UserSalt: "{1,} CaHVhk78uhoPmf5wanA0vg==",
				},
				findSecretErr: errors.New("did not work"),
			},
			false,
			[]assertion{
				func(accountID interface{}) error {
					if cast, ok := accountID.(FindAccountQueryActiveByID); ok {
						if cast != "account-id" {
							return fmt.Errorf("unexpected account identifier %v", cast)
						}
					}
					return nil
				},
				func(evt interface{}) error {
					if cast, ok := evt.(*Event); ok {
						wellformed := cast.Payload == "payload" &&
							cast.AccountID == "account-id" &&
							cast.EventID != "" &&
							cast.SecretID == nil
						if !wellformed {
							return fmt.Errorf("unexpected event shape %v", cast)
						}
					}
					return nil
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := &persistenceLayer{
				dal: test.db,
			}
			err := r.Insert(test.callArgs[0], test.callArgs[1], test.callArgs[2], nil)
			if (err != nil) != test.expectError {
				t.Errorf("Unexpected error value %v", err)
			}
			if expected, found := len(test.argsAssertions), len(test.db.methodArgs); expected != found {
				t.Fatalf("Number of assertions did not match number of calls, got %d and expected %d", found, expected)
			}
			for i, a := range test.argsAssertions {
				if err := a(test.db.methodArgs[i]); err != nil {
					t.Errorf("Unexpected assertion error checking arguments: %v", err)
				}
			}
		})
	}
}

type mockPurgeEventsDatabase struct {
	DataAccessLayer
	findAccountsResult []Account
	findAccountsErr    error
	deleteEventsResult int64
	deleteEventsErr    error
	methodArgs         []interface{}
}

func (m *mockPurgeEventsDatabase) FindAccounts(q interface{}) ([]Account, error) {
	m.methodArgs = append(m.methodArgs, q)
	return m.findAccountsResult, m.findAccountsErr
}

func (m *mockPurgeEventsDatabase) DeleteEvents(q interface{}) (int64, error) {
	m.methodArgs = append(m.methodArgs, q)
	return m.deleteEventsResult, m.deleteEventsErr
}

func (m *mockPurgeEventsDatabase) FindTombstones(q interface{}) ([]Tombstone, error) {
	return nil, nil
}

func (m *mockPurgeEventsDatabase) Commit() error {
	return nil
}

func (m *mockPurgeEventsDatabase) Rollback() error {
	return nil
}

func (m *mockPurgeEventsDatabase) Transaction() (Transaction, error) {
	return m, nil
}

func (m *mockPurgeEventsDatabase) FindEvents(q interface{}) ([]Event, error) {
	return nil, nil
}

func TestPersistenceLayer_Purge(t *testing.T) {
	tests := []struct {
		name          string
		db            *mockPurgeEventsDatabase
		expectError   bool
		argAssertions []assertion
	}{
		{
			"account lookup error",
			&mockPurgeEventsDatabase{
				findAccountsErr: errors.New("did not work"),
			},
			true,
			[]assertion{
				func(q interface{}) error {
					if _, ok := q.(FindAccountsQueryAllAccounts); ok {
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
			},
		},
		{
			"delete events error",
			&mockPurgeEventsDatabase{
				findAccountsResult: []Account{
					{UserSalt: "JF+rNeViJeJb0jth6ZheWg=="},
					{UserSalt: "D6xdWYfRqbuWrkg4OWVgGQ=="},
				},
				deleteEventsErr: errors.New("did not work"),
			},
			true,
			[]assertion{
				func(q interface{}) error {
					if _, ok := q.(FindAccountsQueryAllAccounts); ok {
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
				func(q interface{}) error {
					if hashes, ok := q.(DeleteEventsQueryBySecretIDs); ok {
						for _, hash := range hashes {
							if hash == "user-id" {
								return errors.New("encountered plain user id when hash was expected")
							}
						}
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
			},
		},
		{
			"ok",
			&mockPurgeEventsDatabase{
				findAccountsResult: []Account{
					{UserSalt: "JF+rNeViJeJb0jth6ZheWg=="},
					{UserSalt: "D6xdWYfRqbuWrkg4OWVgGQ=="},
				},
			},
			false,
			[]assertion{
				func(q interface{}) error {
					if _, ok := q.(FindAccountsQueryAllAccounts); ok {
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
				func(q interface{}) error {
					if hashes, ok := q.(DeleteEventsQueryBySecretIDs); ok {
						for _, hash := range hashes {
							if hash == "user-id" {
								return errors.New("encountered plain user id when hash was expected")
							}
						}
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := &persistenceLayer{
				dal: test.db,
			}
			err := r.Purge("user-id")
			if (err != nil) != test.expectError {
				t.Errorf("Unexpected error value %v", err)
			}
			if expected, found := len(test.argAssertions), len(test.db.methodArgs); expected != found {
				t.Fatalf("Number of assertions did not match number of calls, got %d and expected %d", found, expected)
			}
			for i, a := range test.argAssertions {
				if err := a(test.db.methodArgs[i]); err != nil {
					t.Errorf("Assertion error when checking arguments: %v", err)
				}
			}
		})
	}
}

type mockQueryEventDatabase struct {
	DataAccessLayer
	findAccountsResult []Account
	findAccountsErr    error
	findEventsResult   []Event
	findEventsErr      error
	methodArgs         []interface{}
}

func (m *mockQueryEventDatabase) FindAccounts(q interface{}) ([]Account, error) {
	m.methodArgs = append(m.methodArgs, q)
	return m.findAccountsResult, m.findAccountsErr
}

func (m *mockQueryEventDatabase) FindEvents(q interface{}) ([]Event, error) {
	m.methodArgs = append(m.methodArgs, q)
	return m.findEventsResult, m.findEventsErr
}

func (m *mockQueryEventDatabase) FindTombstones(q interface{}) ([]Tombstone, error) {
	return nil, nil
}

func TestPersistenceLayer_Query(t *testing.T) {
	tests := []struct {
		name           string
		db             *mockQueryEventDatabase
		expectedResult EventsResult
		expectError    bool
		argAssertions  []assertion
	}{
		{
			"find accounts error",
			&mockQueryEventDatabase{
				findAccountsErr: errors.New("did not work"),
			},
			EventsResult{},
			true,
			[]assertion{
				func(q interface{}) error {
					if _, ok := q.(FindAccountsQueryAllAccounts); ok {
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
			},
		},
		{
			"find events error",
			&mockQueryEventDatabase{
				findAccountsResult: []Account{
					{UserSalt: "LEWtq55DKObqPK+XEQbnZA=="},
					{UserSalt: "kxwkHp6yPBd0tQ85XlayDg=="},
				},
				findEventsErr: errors.New("did not work"),
			},
			EventsResult{},
			true,
			[]assertion{
				func(q interface{}) error {
					if _, ok := q.(FindAccountsQueryAllAccounts); ok {
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
				func(q interface{}) error {
					if query, ok := q.(FindEventsQueryForSecretIDs); ok {
						if query.Since != "yesterday" {
							return fmt.Errorf("unexpected since value: %v", query.Since)
						}
						if len(query.SecretIDs) != 2 {
							return fmt.Errorf("unexpected number of user ids: %d", len(query.SecretIDs))
						}
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
			},
		},
		{
			"ok",
			&mockQueryEventDatabase{
				findAccountsResult: []Account{
					{AccountID: "account-a", UserSalt: "LEWtq55DKObqPK+XEQbnZA=="},
					{AccountID: "account-b", UserSalt: "kxwkHp6yPBd0tQ85XlayDg=="},
				},
				findEventsResult: []Event{
					{AccountID: "account-a", EventID: "event-a", Payload: "payload-a"},
					{AccountID: "account-b", EventID: "event-b", Payload: "payload-b"},
				},
			},
			EventsResult{
				Events: &EventsByAccountID{
					"account-a": []EventResult{
						{AccountID: "account-a", Payload: "payload-a", EventID: "event-a"},
					},
					"account-b": []EventResult{
						{AccountID: "account-b", Payload: "payload-b", EventID: "event-b"},
					},
				},
			},
			false,
			[]assertion{
				func(q interface{}) error {
					if _, ok := q.(FindAccountsQueryAllAccounts); ok {
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
				func(q interface{}) error {
					if query, ok := q.(FindEventsQueryForSecretIDs); ok {
						if query.Since != "yesterday" {
							return fmt.Errorf("unexpected since value: %v", query.Since)
						}
						if len(query.SecretIDs) != 2 {
							return fmt.Errorf("unexpected number of user ids: %d", len(query.SecretIDs))
						}
						return nil
					}
					return fmt.Errorf("unexpected argument %v", q)
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &persistenceLayer{
				dal: test.db,
			}
			result, err := p.Query(Query{
				UserID: "user-id",
				Since:  "yesterday",
			})

			if (err != nil) != test.expectError {
				t.Errorf("Unexpected error value %v", err)
			}

			if !reflect.DeepEqual(test.expectedResult, result) {
				t.Errorf("Expected %v, got %v", test.expectedResult, result)
			}

			if expected, found := len(test.argAssertions), len(test.db.methodArgs); expected != found {
				t.Fatalf("Number of assertions did not match number of calls, expected %d and found %d", expected, found)
			}

			for i, a := range test.argAssertions {
				if err := a(test.db.methodArgs[i]); err != nil {
					t.Errorf("Assertion error when checking arguments: %v", err)
				}
			}
		})
	}
}

func TestGetLatestSeq(t *testing.T) {
	result := getLatestSeq([]string{"x", "0", "z", "a", "x", "1", "0"})
	if result != "z" {
		t.Errorf("Unexpected result %v", result)
	}
}
