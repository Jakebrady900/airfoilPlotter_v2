package Parsec

import (
	"math"
	"sync"
)

type powerColumnResult struct {
	data       []float64
	identifier float64
}
type multiplyColumnResult struct {
	data       []float64
	identifier string
}

func GenerateAirfoilUpper(coeffMatrix [][]float64) ([]float64, []float64, []float64) {
	Parsec := createParsec(201)

	ch := make(chan powerColumnResult, 6)
	addingCh := make(chan multiplyColumnResult, 12)
	createCoreArrays(Parsec, ch)
	UpperX, LowerX := calculateMainArrays(coeffMatrix, ch, addingCh)

	return UpperX, LowerX, Parsec
}

func powerColumn(indice float64, vector []float64) powerColumnResult {
	result := make([]float64, len(vector))
	for i := 0; i < len(vector); i++ {
		result[i] = math.Pow(vector[i], indice)
	}
	r := powerColumnResult{result, indice}
	return r
}

func multiplyColumn(coefficient float64, vector []float64, identifier string, wg *sync.WaitGroup) multiplyColumnResult {
	defer wg.Done()
	result := make([]float64, len(vector))
	for i := 0; i < len(vector); i++ {
		result[i] = coefficient * vector[i]
	}
	r := multiplyColumnResult{result, identifier}
	return r
}

func createParsec(granularity int) []float64 {
	Parsec := make([]float64, granularity)
	var step float64
	step = 1.0 / float64(granularity-1)
	step = math.Round(step*1e6) / 1e6
	X := 0.0
	for i := 0; i < granularity; i++ {
		Parsec[i] = X
		X += step
	}
	return Parsec
}

func decodeCoefficients(coeffMatrix [][]float64) (a1, a2, a3, a4, a5, a6, b1, b2, b3, b4, b5, b6 float64) {
	a1 = coeffMatrix[0][0]
	a2 = coeffMatrix[1][0]
	a3 = coeffMatrix[2][0]
	a4 = coeffMatrix[3][0]
	a5 = coeffMatrix[4][0]
	a6 = coeffMatrix[5][0]
	b1 = coeffMatrix[6][0]
	b2 = coeffMatrix[7][0]
	b3 = coeffMatrix[8][0]
	b4 = coeffMatrix[9][0]
	b5 = coeffMatrix[10][0]
	b6 = coeffMatrix[11][0]
	return a1, a2, a3, a4, a5, a6, b1, b2, b3, b4, b5, b6

}

func createCoreArrays(Parsec []float64, ch chan powerColumnResult) {

	go func() { ch <- powerColumn(0.5, Parsec) }()
	go func() { ch <- powerColumn(1.5, Parsec) }()
	go func() { ch <- powerColumn(2.5, Parsec) }()
	go func() { ch <- powerColumn(3.5, Parsec) }()
	go func() { ch <- powerColumn(4.5, Parsec) }()
	go func() { ch <- powerColumn(5.5, Parsec) }()
}

func calculateMainArrays(coeffMatrix [][]float64, ch chan powerColumnResult, addingCh chan multiplyColumnResult) (UpperX, LowerX []float64) {
	a1, a2, a3, a4, a5, a6, b1, b2, b3, b4, b5, b6 := decodeCoefficients(coeffMatrix)

	wg := sync.WaitGroup{}
	wg.Add(12)

	for i := 0; i < 6; i++ {
		select {
		case message := <-ch:
			switch message.identifier {
			case 0.5:
				arr_sqrtX := message.data
				go func() { addingCh <- multiplyColumn(a1, arr_sqrtX, "a1", &wg) }()
				go func() { addingCh <- multiplyColumn(b1, arr_sqrtX, "b1", &wg) }()
			case 1.5:
				arr_3rtX := message.data
				go func() { addingCh <- multiplyColumn(a2, arr_3rtX, "a2", &wg) }()
				go func() { addingCh <- multiplyColumn(b2, arr_3rtX, "b2", &wg) }()
			case 2.5:
				arr_5rtX := message.data
				go func() { addingCh <- multiplyColumn(a3, arr_5rtX, "a3", &wg) }()
				go func() { addingCh <- multiplyColumn(b3, arr_5rtX, "b3", &wg) }()
			case 3.5:
				arr_7rtX := message.data
				go func() { addingCh <- multiplyColumn(a4, arr_7rtX, "a4", &wg) }()
				go func() { addingCh <- multiplyColumn(b4, arr_7rtX, "b4", &wg) }()
			case 4.5:
				arr_9rtX := message.data
				go func() { addingCh <- multiplyColumn(a5, arr_9rtX, "a5", &wg) }()
				go func() { addingCh <- multiplyColumn(b5, arr_9rtX, "b5", &wg) }()
			case 5.5:
				arr_11rtX := message.data
				go func() { addingCh <- multiplyColumn(a6, arr_11rtX, "a6", &wg) }()
				go func() { addingCh <- multiplyColumn(b6, arr_11rtX, "b6", &wg) }()
			}
		}
	}
	wg.Wait()
	return listenAndAdd(addingCh)
}

func listenAndAdd(addingCh chan multiplyColumnResult) (UpperX, LowerX []float64) {
	UpperX = make([]float64, 201)
	LowerX = make([]float64, 201)
	for i := 0; i < len(UpperX); i++ {
		UpperX[i] = 0
		LowerX[i] = 0
	}

	mtx := sync.Mutex{}
	for i := 0; i < 12; i++ {
		select {
		case message := <-addingCh:
			switch string(message.identifier[0]) {
			case "a":
				addColumns(UpperX, message.data, &mtx)
			case "b":
				addColumns(LowerX, message.data, &mtx)
			}
		}
	}
	return UpperX, LowerX
}

func addColumns(vector1, vector2 []float64, mtx *sync.Mutex) {
	mtx.Lock()
	defer mtx.Unlock()
	for i := 0; i < len(vector1); i++ {
		vector1[i] += vector2[i]
	}
}
