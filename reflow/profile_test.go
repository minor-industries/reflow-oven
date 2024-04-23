package reflow

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProfile(t *testing.T) {
	err := plot_svg(Profile1, nil, nil)
	require.NoError(t, err)
}
