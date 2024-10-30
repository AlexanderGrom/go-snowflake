// The ID generator is based on the ideas of Twitter Snowflake.
//
// Generates a 64-bit identifier
//  - 41bit contains a timestamp in milliseconds
//  - 10bit contains the number of the machine to generate
//  - 12bit contains the execution number in one millisecond

package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	Timestamp    = 41
	MachineBits  = 10
	SequenceBits = 12
	MaxTimestamp = -1 ^ (-1 << Timestamp)
	MaxMachine   = -1 ^ (-1 << MachineBits)
	MaxSequence  = -1 ^ (-1 << SequenceBits)
)

var (
	ErrMachineID         = errors.New("snowflake: machine id out of range")
	ErrEpochOverflow     = errors.New("snowflake: epoch overflow, unable to generate any more ids")
	ErrTimestampOverflow = errors.New("snowflake: timestamp overflow, unable to generate any more ids")
	ErrInvalidTimestamp  = errors.New("snowflake: invalid timestamp")
)

var timeNow = time.Now

type SnowFlake struct {
	machine        int
	sequence       uint16
	lastTimestamp  int64
	epochTimestamp int64
	lock           sync.Mutex
}

// New creates a new snow flake generator.
func New(machine int, epoch time.Time) (*SnowFlake, error) {
	if machine > MaxMachine {
		return nil, ErrMachineID
	}
	if epoch.After(timeNow()) {
		return nil, ErrEpochOverflow
	}
	return &SnowFlake{
		machine:        machine,
		epochTimestamp: epoch.Unix() * 1e3,
	}, nil
}

// Generate returns a next identity.
func (s *SnowFlake) Generate() (int64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var timestamp = s.timestamp()

	if timestamp > MaxTimestamp {
		return 0, ErrTimestampOverflow
	}

	if timestamp < s.lastTimestamp {
		return 0, ErrInvalidTimestamp
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & MaxSequence
		if s.sequence == 0 {
			timestamp = s.waitNextMillisecond(timestamp)
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return (timestamp << (MachineBits + SequenceBits)) | int64(s.machine<<SequenceBits) | int64(s.sequence), nil
}

func (s *SnowFlake) timestamp() int64 {
	return (timeNow().UnixNano() / 1e6) - s.epochTimestamp
}

func (s *SnowFlake) waitNextMillisecond(last int64) int64 {
	var next = s.timestamp()
	for next <= last {
		next = s.timestamp()
	}
	return next
}
