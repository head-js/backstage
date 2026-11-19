package vault

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"
)

type Record struct {
	Role      string    `json:"role"`
	Secret    string    `json:"secret"`
	N         string    `json:"n"`
	K         string    `json:"k"`
	Timestamp time.Time `json:"timestamp"`
}

type Store struct {
	mu      sync.Mutex
	entries map[string][]Record
}

func New() *Store {
	return &Store{entries: make(map[string][]Record)}
}

func (s *Store) AddSecret() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	for {
		secret := generateSecret()
		if _, exists := s.entries[secret]; !exists {
			s.entries[secret] = nil
			return secret
		}
	}
}

func (s *Store) RemoveKey(secret string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, secret)
}

func (s *Store) Lookup(secret string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.entries[secret]
	return exists
}

func (s *Store) Append(secret, role, n, k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[secret] = append(s.entries[secret], Record{
		Role:      role,
		Secret:    secret,
		N:         n,
		K:         k,
		Timestamp: time.Now(),
	})
}

func (s *Store) All() []Record {
	s.mu.Lock()
	defer s.mu.Unlock()
	var all []Record
	for secret, records := range s.entries {
		for _, r := range records {
			r.Secret = secret
			all = append(all, r)
		}
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Timestamp.After(all[j].Timestamp)
	})
	return all
}

func generateSecret() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		panic("vault: generate secret: " + err.Error())
	}
	return fmt.Sprintf("%06d", n.Int64())
}
