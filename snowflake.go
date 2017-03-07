// Генератор идентификаторов по мотивам Twitter Snowflake.
//
// Генерирует 64-битный идентификатор
//  - 42bit под отметку времени в микросекундах
//  - 10bit под номер машины на котором идет генерация
//  - 12bit под номер выполнения в одну микросекунду

package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	Timestamp    = 42 // этого хватит почти на 140 лет
	MachineBits  = 10
	SequenceBits = 12
	MaxTimestamp = -1 ^ (-1 << Timestamp)
	MaxMachine   = -1 ^ (-1 << MachineBits)
	MaxSequence  = -1 ^ (-1 << SequenceBits)
)

// Timestamp начала эпохи генерации индетификаторов
// Установите свой, например, дату начала разработки проекта
var epochTimestamp int64 = time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix() * 1e3

var (
	ErrMachineID         = errors.New("Snowflake: Machine ID out of range")
	ErrTimestampOverflow = errors.New("Snowflake: TimeStamp overflow. Unable to generate any more IDs")
	ErrInvalidTimestamp  = errors.New("Snowflake: Invalid timestamp")
)

type SnowFlake struct {
	machine       uint16
	sequence      uint16
	lastTimestamp uint64
	lock          sync.Mutex
}

func New(machine uint16) (*SnowFlake, error) {
	if machine > MaxMachine {
		return nil, ErrMachineID
	}
	return &SnowFlake{
		machine: machine,
	}, nil
}

func (self *SnowFlake) Generate() (uint64, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	timestamp := timestamp()

	if timestamp > MaxTimestamp {
		return 0, ErrTimestampOverflow
	}

	if timestamp < self.lastTimestamp {
		return 0, ErrInvalidTimestamp
	}

	if timestamp == self.lastTimestamp {
		self.sequence = (self.sequence + 1) & MaxSequence
		if self.sequence == 0 {
			timestamp = waitNextMicrosecond(timestamp)
		}
	} else {
		self.sequence = 0
	}

	self.lastTimestamp = timestamp

	return uint64((timestamp << (MachineBits + SequenceBits)) | uint64(self.machine<<MachineBits) | uint64(self.sequence)), nil
}

// Date format YYYY-MM-DD
func EpochTimestamp(date string) error {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		epochTimestamp = t.Unix() * 1e3
	}
	return err
}

func timestamp() uint64 {
	return uint64((time.Now().UnixNano() / 1e6) - epochTimestamp)
}

func waitNextMicrosecond(last uint64) uint64 {
	ts := timestamp()
	for ts < last {
		ts = timestamp()
	}
	return ts
}
