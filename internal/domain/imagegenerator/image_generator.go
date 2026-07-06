package imagegenerator

import (
	"crypto/rand"
	"encoding/binary"
	"image/color"
	"math"
	"math/big"
	"sync"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/function"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/constants"
)

// Описывает точку изображения.
type Point struct {
	X       int
	Y       int
	Color   color.NRGBA
	counter int     // Число попаданий в точку.
	normal  float64 // Используется для коррекции изображения
}

// Описывает аффинное преобразование.
type affine struct {
	a     float64
	b     float64
	c     float64
	d     float64
	e     float64
	f     float64
	red   uint8
	green uint8
	blue  uint8
}

// Описывает генератор изображения.
type Generator struct {
	wg *sync.WaitGroup
}

// Конструктор генератора изображения.
func NewGenerator(wg *sync.WaitGroup) *Generator {
	return &Generator{wg: wg}
}

// Инициализирует необходимые структуры данных, вызывает горутины генерации точек.
func (gen *Generator) Generate(pointsCount, affinesCount, it, xRes, yRes int, isSymmetry bool, funcs []function.NonLinearFunc) *[][]Point {
	// Генерация коэффициентов аффинных преобразований.
	affines := make([]affine, affinesCount)

	for i := 0; i < affinesCount; {
		a := getRandomFloat64(-1, 1)
		b := getRandomFloat64(-1, 1)
		c := getRandomFloat64(-1, 1)
		d := getRandomFloat64(-1, 1)
		e := getRandomFloat64(-3, 3)
		f := getRandomFloat64(-3, 3)

		if (a*a+d*d < 1) && (b*b+e*e < 1) && (a*a+b*b+d*d+e*e < 1+(a*e-b*d)*(a*e-b*d)) {
			red := getRandomUint8()
			green := getRandomUint8()
			blue := getRandomUint8()
			affines[i] = affine{a, b, c, d, e, f, red, green, blue}
			i++
		}
	}

	points := make([][]Point, yRes)
	for i := 0; i < yRes; i++ {
		points[i] = make([]Point, xRes)
	}

	for i := 0; i < yRes; i++ {
		for j := 0; j < xRes; j++ {
			points[i][j].Y = i
			points[i][j].X = j
			points[i][j].Color = color.NRGBA{0, 0, 0, 255}
		}
	}

	pointsPerGoroutine := pointsCount / constants.GoroutineCount

	for i := 0; i < constants.GoroutineCount; i++ {
		gen.wg.Add(1)
		go gen.process(pointsPerGoroutine, it, affinesCount, xRes, yRes, affines, funcs, points, isSymmetry)
	}

	return &points
}

// Генерирует точки изображения.
func (gen *Generator) process(pointsCount, it, affinesCount, xRes, yRes int, affines []affine,
	funcs []function.NonLinearFunc, points [][]Point, isSymmetry bool) {
	defer func() {
		gen.wg.Done()
	}()

	for i := 0; i < pointsCount; i++ {
		const (
			XMIN = -3.54
			YMIN = -2
			YMAX = 2
			XMAX = 3.54
		)

		newX := getRandomFloat64(XMIN, XMAX)
		newY := getRandomFloat64(YMIN, YMAX)

		// Первые 20 итераций точка не рисуется.
		for j := -20; j < it; j++ {
			randInt, _ := getRandomInt(0, affinesCount)

			// Применяются аффинные преобразования.
			x := affines[randInt].a*newX + affines[randInt].b*newY + affines[randInt].c
			y := affines[randInt].d*newX + affines[randInt].e*newY + affines[randInt].f

			// Применяются variants - нелинейные функции.
			x, y = applyFunctions(x, y, funcs)

			if j >= 0 && (x >= XMIN && x <= XMAX) && (y >= YMIN && y <= YMAX) {
				x1 := xRes - int(math.Trunc(((XMAX-x)/(XMAX-XMIN))*float64(xRes)))
				y1 := yRes - int(math.Trunc(((YMAX-y)/(YMAX-YMIN))*float64(yRes)))

				// Слайс с точками x, при симметрии необходимо хранить координаты 2ух точек - левой и правой.
				var xs []int

				center := xRes / 2

				switch {
				case isSymmetry && x1 >= center:
					xs = []int{x1, center - (x1 - center)} // x1 - center - расстояние от центра до 1ой точки.
				case isSymmetry && x1 < center:
					xs = []int{x1, center + (center - x1)}
				default:
					xs = []int{x1}
				}

				// Обходятся все x: 1 если нет симметрии, 2 если есть симметрия.
				for _, x := range xs {
					// Если точка попадает в область рисунка, то рисуем ее.
					if x < xRes && y1 < yRes {
						colorPoint(x, y1, points, affines[randInt])
					}
				}
			}
		}
	}
}

// Раскрашивает точку.
func colorPoint(x1, y1 int, points [][]Point, affineTransform affine) {
	if points[y1][x1].counter == 0 {
		points[y1][x1].Color.R = affineTransform.red
		points[y1][x1].Color.G = affineTransform.green
		points[y1][x1].Color.B = affineTransform.blue
	}

	alpha := getRandomUint8()
	points[y1][x1].Color.A = alpha
	points[y1][x1].counter++
}

// Применяет к точке случайную функцию из списка.
func applyFunctions(x, y float64, funcs []function.NonLinearFunc) (newX, newY float64) {
	variantionIndex, _ := getRandomInt(0, len(funcs))
	newX, newY = funcs[variantionIndex](x, y)

	return newX, newY
}

// Возвращает случайный int.
func getRandomInt(lowerBound, upperBound int) (int, error) {
	diff := big.NewInt(int64(upperBound - lowerBound))

	randomValue, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return -1, GeneratorError{"ошибка генерации случайного int", err}
	}

	return int(randomValue.Int64()) + lowerBound, nil
}

// Возвращает случайный float64.
func getRandomFloat64(lowerBound, upperBound float64) float64 {
	var randomInt uint64
	_ = binary.Read(rand.Reader, binary.LittleEndian, &randomInt)

	randomFraction := float64(randomInt) / (1 << 64)

	return lowerBound + randomFraction*(upperBound-lowerBound)
}

// Возвращает случайный uint8.
func getRandomUint8() uint8 {
	var randomInt uint64
	_ = binary.Read(rand.Reader, binary.LittleEndian, &randomInt)

	// Случайное число от 0 до 1
	randomFraction := float64(randomInt) / (1 << 64)

	lowerBound := uint8(1)
	upperBound := uint8(255)

	return lowerBound + uint8(randomFraction*float64(upperBound-lowerBound))
}

// Применяет к точкам гамма коррекцию для избавления от шума.
func Correct(width, height int, points *[][]Point) {
	maxValue := 0.0
	gamma := 2.2

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if (*points)[i][j].counter != 0 {
				(*points)[i][j].normal = math.Log10(float64((*points)[i][j].counter))
				maxValue = math.Max(maxValue, (*points)[i][j].normal)
			}
		}
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			(*points)[i][j].normal /= maxValue
			(*points)[i][j].Color.R = uint8(float64((*points)[i][j].Color.R) * (math.Pow((*points)[i][j].normal, (1.0 / gamma))))
			(*points)[i][j].Color.G = uint8(float64((*points)[i][j].Color.G) * (math.Pow((*points)[i][j].normal, (1.0 / gamma))))
			(*points)[i][j].Color.B = uint8(float64((*points)[i][j].Color.B) * (math.Pow((*points)[i][j].normal, (1.0 / gamma))))
		}
	}
}
