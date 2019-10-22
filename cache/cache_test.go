package cache

import "testing"

func TestLRUSize(t *testing.T) {
	KEYONE := "name"
	KEYTWO := "place"
	VALUEONE := "Andrew"
	VALUETWO := "Moon"

	// Test a cache with keys limited to one
	s := Service(1)

	s.Set(KEYONE, VALUEONE, 0)

	if curValue, ok := s.Get(KEYONE); ok != true || curValue != VALUEONE {
		t.Errorf("Problem setting and getting key")
	}

	s.Set(KEYTWO, VALUETWO, 0)

	if _, ok := s.Get(KEYONE); ok == true {
		t.Errorf("First key wasn't removed from cache as it become oversized")
	}
}

func TestLRUExpiring(t *testing.T) {
	KEYONE := "name"
	VALUEONE := "Andrew"
	EXPIRE := 1

	// Test a cache with an expired key
	s := Service(1)

	s.Set(KEYONE, VALUEONE, EXPIRE)

	if _, ok := s.Get(KEYONE); ok == true {
		t.Errorf("Key should have expired")
	}
}
