package utils

import "strconv"

// ConvertPort converts a port number given as a string to an unsigned 16-bit integer
// All unix ports are within the range of 0-65565, which is exactly the range of an unsigned 16-bit integer
func ConvertPort(p string) (uint16, error) {
	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}

	return uint16(i), nil
}
