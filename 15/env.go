package main

import "os"

func expandEnv(input string) string {
	return os.Expand(input, func(key string) string {
		return os.Getenv(key)
	})
}
