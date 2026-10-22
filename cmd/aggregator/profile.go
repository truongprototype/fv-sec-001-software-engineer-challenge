package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/pprof"
)

func withCPUProfile(path string, run func() error) (runErr error) {
	if path == "" {
		return run()
	}

	profileFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create CPU profile %q: %w", path, err)
	}
	if err := pprof.StartCPUProfile(profileFile); err != nil {
		closeErr := profileFile.Close()
		profileErr := fmt.Errorf("start CPU profile %q: %w", path, err)
		if closeErr != nil {
			return errors.Join(profileErr, fmt.Errorf("close CPU profile %q: %w", path, closeErr))
		}
		return profileErr
	}

	defer func() {
		pprof.StopCPUProfile()
		if err := profileFile.Sync(); err != nil {
			runErr = errors.Join(runErr, fmt.Errorf("sync CPU profile %q: %w", path, err))
		}
		if err := profileFile.Close(); err != nil {
			runErr = errors.Join(runErr, fmt.Errorf("close CPU profile %q: %w", path, err))
		}
	}()

	return run()
}
