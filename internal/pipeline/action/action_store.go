package action

import (
	"errors"
	"fmt"
)

type ActionStore struct {
	actions map[string]*Action
}

func NewActionStore() *ActionStore {
	as := &ActionStore{}
	as.actions = make(map[string]*Action)
	return as
}

func (s *ActionStore) Get(key string) (*Action, error) {
	_action, exist := s.actions[key]
	if exist {
		return _action, nil
	}

	return nil, errors.New(fmt.Sprintf("Action not found with key: %s", key))
}

func (s *ActionStore) Put(key string, _action Action) {
	s.actions[key] = &_action
}