package ui_test

import (
	"io"
	"os"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/ui"

	"github.com/stretchr/testify/assert"
)

func Test_CheckNumbers(t *testing.T) {
	type TestCase struct {
		name       string
		slice      []int
		lowerBound int
		upperBound int
		expected   bool
	}

	TestCases := []TestCase{
		{name: "Все значение - корректные", slice: []int{3, 4, 5, 1}, lowerBound: 1, upperBound: 6, expected: true},
		{name: "Одно значение - некорректное", slice: []int{3, 10, 5, 1}, lowerBound: 1, upperBound: 6, expected: false},
		{name: "Все значения - одинаковые коррекные", slice: []int{3, 3, 3, 3}, lowerBound: 1, upperBound: 6, expected: true},
	}

	UI := ui.NewUI(os.Stdin, io.Discard)
	for _, tc := range TestCases {
		assert.Equal(t, tc.expected, UI.CheckNumbers(tc.slice, tc.lowerBound, tc.upperBound), tc.name)
	}
}

func Test_ProcessNumbers(t *testing.T) {
	type TestCase struct {
		name     string
		input    string
		expected []int
	}

	TestCases := []TestCase{
		{name: "Корректный ввод (несколько номеров)", input: "1 2 3 4", expected: []int{0, 1, 2, 3}},
		{name: "Корректный ввод (один номер)", input: "1", expected: []int{0}},
		{name: "Некорректный ввод (неправильное кол-во пробелов)", input: " 1 3  4 5", expected: nil},
		{name: "Некорректный ввод (символы - не числа)", input: ", % $ @", expected: nil},
	}

	UI := ui.NewUI(os.Stdin, io.Discard)
	for _, tc := range TestCases {
		// В функции есть только одно место, которое может вернуть err - Atoi. Если она возвращает error, то возвращается nil.
		res, _ := UI.ProcessNumbers(tc.input)
		assert.Equal(t, tc.expected, res, tc.name)
	}
}
