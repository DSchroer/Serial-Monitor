package main

import (
	"errors"
)

type monitor struct {
	name   string
	device string
}

func (m monitor) readIntArray(message []byte, searchId []byte, length int) ([]int, error) {
	data, err := requestData(message, m.device)
	if err != nil {
		return nil, err
	}

	start := idPos(searchId, data)

	if start == -1 {
		return nil, errors.New("Data not found")
	}

	values := make([]int, length)
	for i := 0; i < length; i++ {
		values[i] = readLEInt(data, start+(i*2))
	}

	return values, nil
}

func (m monitor) currentDraw() (int, error) {
	current := []byte{0xAA, 0x02, 0x00, 0xAD}
	wattId := []byte{0x53, 0x30, 0x32}

    data, err := m.readIntArray(current, wattId, 2)
    if err != nil {
        return 0,nil
    }

	return data[1], nil
}

func (m monitor) dailyUsage() ([]int, error) {
	day := []byte{0xAA, 0x01, 0x00, 0xAD}
	dayId := []byte{0x53, 0x30, 0x31, 0x44, 0x41, 0x59}

	return m.readIntArray(day, dayId, 8)
}

func (m monitor) weeklyUsage() ([]int, error) {
	week := []byte{0x0A, 0x0A, 0x00, 0xAD}
	weekId := []byte{0x53, 0x30, 0x31, 0x57, 0x45, 0x4B}

	return m.readIntArray(week, weekId, 8)
}
