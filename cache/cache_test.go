package cache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lulzshadowwalker/pokedex/cache"
)

func TestRememberAndGet(t *testing.T) {
	const interval = 5 * time.Second
	const duration = 5 * time.Second

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := cache.New(interval)
			defer cache.Stop()

			resolver := func() []byte {
				return c.val
			}

			cache.Remember(c.key, resolver, duration)

			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}

			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestCleanupLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 10*time.Millisecond

	cache := cache.New(baseTime)
	defer cache.Stop()

	resolver := func() []byte {
		return []byte("testdata")
	}

	cache.Remember("https://example.com", resolver, baseTime)

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestForget(t *testing.T) {
	const interval = 5 * time.Second

	cache := cache.New(interval)
	defer cache.Stop()

	resolver := func() []byte {
		return []byte("testdata")
	}

	cache.Remember("https://example.com", resolver, interval)
	cache.Forget("https://example.com")

	_, ok := cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
	}
}
