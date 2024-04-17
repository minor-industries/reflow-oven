package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProfile(t *testing.T) {
	err := plot_svg(profile1, nil, nil)
	require.NoError(t, err)
}
