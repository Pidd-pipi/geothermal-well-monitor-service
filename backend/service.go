package main

type WellService struct{ store *WellStore }

func NewWellService(store *WellStore) *WellService { return &WellService{store: store} }
func (s *WellService) Wells() []Well               { return s.store.List() }
func (s *WellService) ChangeStatus(id, status string) (Well, error) {
	if err := ValidateWellStatus(status); err != nil {
		return Well{}, err
	}
	return s.store.UpdateStatus(id, status)
}
