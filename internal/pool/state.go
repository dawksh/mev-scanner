package pool

import (
	"math/big"
	"sync"
)

type State struct {
	mu       sync.RWMutex
	Reserve0 *big.Int
	Reserve1 *big.Int
}

func NewState() *State {
	return &State{
		Reserve0: big.NewInt(0),
		Reserve1: big.NewInt(0),
	}
}

func (s *State) Update(r0, r1 *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Reserve0 = new(big.Int).Set(r0)
	s.Reserve1 = new(big.Int).Set(r1)
}

func (s *State) Get() (*big.Int, *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return new(big.Int).Set(s.Reserve0), new(big.Int).Set(s.Reserve1)
}
