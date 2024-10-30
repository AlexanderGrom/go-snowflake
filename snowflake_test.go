package snowflake

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenID(t *testing.T) {
	epoch := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)
	timeNow = time.Date(2016, 1, 1, 1, 0, 0, 0, time.UTC).UTC

	gen, err := New(1, epoch)
	require.NoError(t, err)

	num, err := gen.Generate()
	require.NoError(t, err)
	require.Equal(t, int64(15099494404096), num)

	timeNow = time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC).UTC

	num, err = gen.Generate()
	require.NoError(t, err)
	require.Equal(t, int64(132633958809604096), num)

	timeNow = time.Date(2080, 1, 1, 0, 0, 0, 0, time.UTC).UTC

	num, err = gen.Generate()
	require.NoError(t, err)
	require.Equal(t, int64(8471178746265604096), num)
}

func TestMaxMachine(t *testing.T) {
	epoch := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)
	timeNow = time.Date(2016, 1, 1, 1, 0, 0, 0, time.UTC).UTC

	var err error
	_, err = New(1, epoch)
	require.NoError(t, err)

	_, err = New(1023, epoch)
	require.NoError(t, err)

	_, err = New(1024, epoch)
	require.ErrorIs(t, err, ErrMachineID)
}

func TestBadEpoch(t *testing.T) {
	epoch := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	timeNow = time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).UTC

	_, err := New(1, epoch)
	require.ErrorIs(t, err, ErrEpochOverflow)
}

func TestTimeOverflow(t *testing.T) {
	epoch := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

	gen, err := New(1, epoch)
	require.NoError(t, err)

	timeNow = time.Date(2085, 9, 6, 0, 0, 0, 0, time.UTC).UTC

	_, err = gen.Generate()
	require.NoError(t, err)

	timeNow = time.Date(2085, 9, 7, 0, 0, 0, 0, time.UTC).UTC

	_, err = gen.Generate()
	require.ErrorIs(t, err, ErrTimestampOverflow)
}

func TestSequence(t *testing.T) {
	epoch := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

	gen, err := New(1, epoch)
	require.NoError(t, err)

	sequence := 0
	timeNow = func() time.Time {
		sequence++
		if sequence <= 4100 {
			return time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC).UTC()
		}
		return time.Date(2017, 1, 1, 0, 0, 1, 0, time.UTC).UTC()
	}

	for i := 0; i <= 4095; i++ {
		num, err := gen.Generate()
		require.NoError(t, err)
		require.Equal(t, int64(132633958809604096+i), num)
	}

	num, err := gen.Generate()
	require.NoError(t, err)
	require.Equal(t, int64(132633963003908096), num)
}
