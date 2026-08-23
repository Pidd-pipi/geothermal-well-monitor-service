package main

import "fmt"

type WellService struct{ store *WellStore }

func NewWellService(store *WellStore) *WellService { return &WellService{store: store} }
func (s *WellService) Wells() []Well               { return s.store.List() }
func (s *WellService) ChangeStatus(id, status string) (Well, error) {
	well, err := s.store.UpdateStatus(id, status)
	if err != nil {
		return Well{}, fmt.Errorf("update well %s: %v", id, err)
	}
	if err := ValidateWellStatus(status); err != nil {
		return Well{}, err
	}
	return well, nil
}
