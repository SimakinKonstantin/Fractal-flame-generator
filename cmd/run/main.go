package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/imagegenerator"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/constants"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/ui"
)

func main() {
	UI := ui.NewUI(os.Stdin, os.Stdout)

	width, height, err := UI.InputResolution()
	if err != nil {
		fmt.Println("ошибка работы приложения", err)
	}

	iterationCount, err := UI.InputIterationsCount()
	if err != nil {
		fmt.Println("ошибка работы приложения", err)
	}

	funcs, err := UI.InputTransformFunctions()
	if err != nil {
		fmt.Println("ошибка работы приложения", err)
	}

	isSymmetry, err := UI.InputSymmetry()
	if err != nil {
		fmt.Println("ошибка работы приложения", err)
	}

	err = UI.PrintWaitMsg()
	if err != nil {
		fmt.Println("ошибка работы приложения", err)
	}

	start := time.Now()
	wg := sync.WaitGroup{}
	generator := imagegenerator.NewGenerator(&wg)

	imageScheme := generator.Generate(constants.PointsCount, constants.AffinesCount, iterationCount, width, height, isSymmetry, funcs)

	wg.Wait()

	imagegenerator.Correct(width, height, imageScheme)

	img := UI.CreateImage(imageScheme)

	err = UI.PrintImage(img)
	if err != nil {
		fmt.Println("ошибка работы приложения", err)
	}

	duration := time.Since(start)
	fmt.Println(duration)
}
