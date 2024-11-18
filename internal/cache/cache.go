package cache

import (
	"time"
	"nerd-fonts-cli/pkg/cache"
	"nerd-fonts-cli/pkg/data"
)

// Refresh refreshes the cache.
func Refresh(cache cache.Cache) error {
	resp, err := data.Fetch()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return cache.Write(resp.Body)
}

// RefreshIfOld refreshes the cache if it is too old.
func RefreshIfOld(cache cache.Cache, maxAge time.Duration) error {
	// NOTE We ignore the error, assuming the cache doesn't exist if there was an issue
	//		checking its age.
	age, _ := cache.Age()
	if age > maxAge {
		return Refresh(cache)
	}
	return nil
}
