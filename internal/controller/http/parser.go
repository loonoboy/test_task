package http

import (
	"fmt"
	"strconv"
	"strings"
)

func parseIDFromURL(path, prefix string) (int, error) {
	if !strings.HasPrefix(path, prefix) {
		return 0, fmt.Errorf("invalid path")
	}

	idStr := strings.TrimPrefix(path, prefix)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return id, nil
}
