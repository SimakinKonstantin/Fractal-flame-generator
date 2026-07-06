package function

import "math"

// Описывает тип нелинейных функций.
type NonLinearFunc = func(float64, float64) (float64, float64)

func Handkerchief(x, y float64) (newX, newY float64) {
	newX = math.Sqrt(x*x+y*y) * math.Sin(x/y+math.Sqrt(x*x+y*y))
	newY = math.Sqrt(x*x+y*y) * math.Cos(x/y-math.Sqrt(x*x+y*y))

	return newX, newY
}

func Disc(x, y float64) (newX, newY float64) {
	newX = math.Atan(x/y) / math.Pi * math.Sin(math.Pi*math.Sqrt(x*x+y*y))
	newY = math.Atan(x/y) / math.Pi * math.Cos(math.Pi*math.Sqrt(x*x+y*y))

	return newX, newY
}

func Heart(x, y float64) (newX, newY float64) {
	newX = math.Sqrt(x*x+y*y) * math.Sin(math.Sqrt(x*x+y*y)*math.Atan(y/x))
	newY = -math.Sqrt(x*x+y*y) * math.Cos(math.Sqrt(x*x+y*y)*math.Atan(y/x))

	return newX, newY
}

func Diamond(x, y float64) (newX, newY float64) {
	newX = math.Sin(math.Atan(x/y)) * math.Cos(math.Sqrt(x*x+y*y))
	newY = math.Cos(math.Atan(x/y)) * math.Sin(math.Sqrt(x*x+y*y))

	return newX, newY
}

func Swirl(x, y float64) (newX, newY float64) {
	newX = x*math.Sin(x*x+y*y) - y*math.Cos(x*x+y*y)
	newY = x*math.Cos(x*x+y*y) + y*math.Sin(x*x+y*y)

	return newX, newY
}
