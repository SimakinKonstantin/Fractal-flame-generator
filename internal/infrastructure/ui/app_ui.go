package ui

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/function"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/imagegenerator"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/constants"
)

// Описывает пользовательский интерфейс.
type UI struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

// Создает объект пользовательского интерфейса.
func NewUI(reader io.Reader, writer io.Writer) *UI {
	return &UI{reader: bufio.NewReader(reader), writer: bufio.NewWriter(writer)}
}

// Используется для ввода разрешения изображения.
func (ui UI) InputResolution() (width, height int, err error) {
	for {
		_, err := fmt.Fprint(ui.writer, "Введите ширину изображения в пикселях: ")
		if err != nil {
			return -1, -1, OutputError{"ошибка вывода запроса ширины изображения", err}
		}

		err = ui.writer.Flush()
		if err != nil {
			return -1, -1, OutputError{"ошибка вывода запроса ширины изображения", err}
		}

		widthStr, _, err := ui.reader.ReadLine()
		if err != nil {
			return -1, -1, InputError{"ошибка ввода ширины изображения", err}
		}

		width, err = strconv.Atoi(string(widthStr))
		if err != nil || width <= 0 {
			continue
		}

		break
	}

	for {
		_, err = fmt.Fprint(ui.writer, "Введите высоту изображения в пикселях: ")
		if err != nil {
			return -1, -1, OutputError{"ошибка вывода запроса высоты изображения", err}
		}

		err = ui.writer.Flush()
		if err != nil {
			return -1, -1, OutputError{"ошибка вывода запроса высоты изображения", err}
		}

		heightStr, _, err := ui.reader.ReadLine()
		if err != nil {
			return -1, -1, InputError{"ошибка ввода высоты изображения", err}
		}

		height, err = strconv.Atoi(string(heightStr))
		if err != nil || height <= 0 {
			continue
		}

		return width, height, nil
	}
}

// Используется для ввода кол-ва итераций.
func (ui UI) InputIterationsCount() (count int, err error) {
	for {
		_, err := fmt.Fprint(ui.writer, "Введите количество итераций: ")
		if err != nil {
			return -1, OutputError{"ошибка вывода запроса количества итераций", err}
		}

		err = ui.writer.Flush()
		if err != nil {
			return -1, OutputError{"ошибка вывода запроса количества итераций", err}
		}

		countStr, _, err := ui.reader.ReadLine()
		if err != nil {
			return -1, InputError{"ошибка ввода при запросе количества итераций", err}
		}

		count, err = strconv.Atoi(string(countStr))
		if err != nil || count < 1 {
			continue
		}

		return count, nil
	}
}

// Используется для ввода информации о симметричности изображения.
func (ui UI) InputSymmetry() (isSymmetry bool, err error) {
	for {
		_, err = fmt.Fprint(ui.writer, "Введите 'Y' - изображение должно быть симметричным, иначе - 'N': ")
		if err != nil {
			return false, OutputError{"ошибка вывода запроса информации о симметрии", err}
		}

		err = ui.writer.Flush()
		if err != nil {
			return false, OutputError{"ошибка вывода запроса информации о симметрии", err}
		}

		isSymmetryBytes, _, err := ui.reader.ReadLine()
		if err != nil {
			return false, InputError{"ошибка ввода информации о симметрии", err}
		}

		// Проверка корректности ввода.
		isSymmetryStr := strings.ToUpper(string(isSymmetryBytes))
		switch isSymmetryStr {
		case "Y":
			return true, nil
		case "N":
			return false, nil
		default:
			continue
		}
	}
}

// Используется для выбор нелинейных функций, которые будут использоваться при генерации.
func (ui UI) InputTransformFunctions() ([]function.NonLinearFunc, error) {
	for {
		_, err := fmt.Fprint(ui.writer, "Введите номера функций через пробел (формат: [число1 число2 число3]): \n")
		if err != nil {
			return nil, OutputError{"ошибка вывода запроса номеров функций", err}
		}

		funcs := []string{
			"Handkerchief",
			"Disc",
			"Heart",
			"Diamond",
			"Swirl",
		}

		for ind, funcName := range funcs {
			_, err := fmt.Fprintf(ui.writer, "%d. %s\n", ind+1, funcName)
			if err != nil {
				return nil, OutputError{"ошибка вывода списка доступных функций", err}
			}
		}

		err = ui.writer.Flush()
		if err != nil {
			return nil, OutputError{"ошибка вывода", err}
		}

		NumbersBytes, _, err := ui.reader.ReadLine()
		if err != nil {
			return nil, InputError{"ошибка ввода номера функции", err}
		}

		// Если не удалось получить ожидаемый формат данных, ввод запрашивается снова.
		NumbersSlice, err := ui.ProcessNumbers(string(NumbersBytes))
		if err != nil {
			continue
		}

		// min = 0, т.к. в ProcessNumbers введенный номер приводится к формату индекс массива (1->0, 2->1 ...).
		if !ui.CheckNumbers(NumbersSlice, 0, len(funcs)) {
			continue
		}

		return ui.getFuncsByNumbers(NumbersSlice, funcs), nil
	}
}

// Обрабатывает номера функций, который выбрал пользователь.
// Все номера, который ввел пользователь уменьшает на 1, чтобы далее можно было работать с индексами слайсов.
func (ui UI) ProcessNumbers(input string) ([]int, error) {
	numbersStr := strings.Split(input, " ")

	numbersInt := make([]int, len(numbersStr))

	var err error
	for ind, numberStr := range numbersStr {
		numbersInt[ind], err = strconv.Atoi(numberStr)
		if err != nil {
			return nil, err
		}

		// Т.к. выводимый ui список функций нумеруется с 1.
		numbersInt[ind]--
	}

	return numbersInt, err
}

// Проверяет вхождение выбранных пользователем номеров в корректный диапазон.
func (ui UI) CheckNumbers(nums []int, lowerBound, upperBound int) bool {
	for _, num := range nums {
		if num > upperBound || num < lowerBound {
			return false
		}
	}

	return true
}

// Сопоставляет выбранные пользователем номера функций с самими функциями.
func (ui UI) getFuncsByNumbers(numbers []int, funcs []string) []function.NonLinearFunc {
	result := make([]function.NonLinearFunc, len(numbers))

	for ind, number := range numbers {
		fun := funcs[number]
		switch fun {
		case "Handkerchief":
			result[ind] = function.Handkerchief
		case "Disc":
			result[ind] = function.Disc
		case "Heart":
			result[ind] = function.Heart
		case "Diamond":
			result[ind] = function.Diamond
		case "Swirl":
			result[ind] = function.Swirl
		}
	}

	return result
}

// Печатает сообщение о старте работы программы.
func (ui UI) PrintWaitMsg() error {
	_, err := fmt.Fprint(ui.writer, "Генерация запущена. Пожалуйста подождите...\n")
	if err != nil {
		return OutputError{"ошибка вывода сообщения об ожидании", err}
	}

	err = ui.writer.Flush()
	if err != nil {
		return OutputError{"ошибка вывода сообщения об ожидании", err}
	}

	return nil
}

// Создает объект изображения.
func (ui UI) CreateImage(points *[][]imagegenerator.Point) *image.NRGBA {
	height := len(*points)
	width := len((*points)[0])

	// Исходя из того, что points - матрица точек, определяются размеры изображения.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set((*points)[y][x].X, (*points)[y][x].Y, (*points)[y][x].Color)
		}
	}

	return img
}

// todo defer FileClose
// Записывает объект изображения в файл, путь к которому указан в пакете constants.
func (ui UI) PrintImage(img *image.NRGBA) (err error) {
	file, err := os.Create(constants.FileName)
	defer func() {
		// Закрытия файла с перезаписью err, если при закрытии произошла ошибка.
		if CloseErr := file.Close(); CloseErr != nil {
			err = FileError{"ошибка закрытия файла", CloseErr}
		}
	}()

	if err != nil {
		return FileError{"ошибка создания файла", err}
	}

	if err = png.Encode(file, img); err != nil {
		return FileError{"ошибка записи файла", err}
	}

	return nil
}
