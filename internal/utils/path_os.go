package utils

import "os"

func osUserHomeDir() (string, error) { return os.UserHomeDir() }
