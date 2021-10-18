package models

import (
	"database/sql/driver"
	"errors"
	"strconv"
)

type Path struct {
	Points []Point
}

func (p Path) GormDataType() string {
	return "path"
}

func (p Path) Value() (driver.Value, error) {
	out := []byte{'('}
	for _, point := range p.Points {
		v, _ := point.Value()
		switch value := v.(type) {
		case string:
			src := []byte(value)
			out = append(out, src...)
		case []byte:
			out = append(out, value...)
		}

		out = append(out, ',')
	}
	out = out[:len(out)-1]
	out = append(out, ')')

	return string(out), nil
}

func (p *Path) Scan(src interface{}) (err error) {
	var data []byte
	switch src := src.(type) {
	case []byte:
		data = src
	case string:
		data = []byte(src)
	case nil:
		return nil
	default:
		return errors.New("(*Point).Scan: unsupported data type")
	}

	if len(data) == 0 {
		return nil
	}

	for i := len(data) - 1; i >= 0; i-- {
		switch data[i] {
		case '(':
			data = remove(data, i)
		case ')':
			data = remove(data, i)
		}
	}

	var numbers [][]byte
	lastIndex := 0

	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			numbers = append(numbers, data[lastIndex:i])
			lastIndex = i + 1
		}
	}

	numbers = append(numbers, data[lastIndex:])

	for i := 0; i < len(numbers); i += 2 {
		x, _ := strconv.ParseFloat(string(numbers[i]), 64)
		y, _ := strconv.ParseFloat(string(numbers[i+1]), 64)
		p.Points = append(p.Points, Point{
			X: x,
			Y: y,
		})
	}

	return nil
}

func remove(slice []byte, index int) []byte {
	return append(slice[:index], slice[index+1:]...)
}
