package P_Matrix

import (
	"math"
)

type RowResult struct {
	result []float64
	index  int
}

func CreateP(R_LE, Yu, d2Yu, Yl, d2Yl, del_Y_TE, Y_TE, a_TE, b_TE float64) [][]float64 {
	a1 := math.Pow((2 * R_LE), 0.5)
	a2 := Yu
	a3 := d2Yu
	a4 := Yl
	a5 := d2Yl
	a6 := del_Y_TE
	b1 := Y_TE
	b2 := math.Tan(a_TE*(math.Pi/180)) - math.Tan((b_TE/2)*(math.Pi/180))
	b3 := math.Tan(a_TE*(math.Pi/180)) + math.Tan((b_TE/2)*(math.Pi/180))
	b4 := 0.00
	b5 := 0.00
	b6 := -1 * (math.Pow(2*R_LE, 0.5))

	values := []float64{a1, a2, a3, a4, a5, a6, b1, b2, b3, b4, b5, b6}
	P := make([][]float64, 12)
	for x := 0; x < 12; x++ {
		P[x] = []float64{0}
		P[x][0] = values[x]
	}

	return P

}

func Multiply(mat1 [][]float64, mat2 [][]float64) [][]float64 {

	LenA := len(mat1)
	colB := len(mat2[0])

	finalResult := make([][]float64, LenA)
	rowChannel := make(chan RowResult, LenA)

	for i := 0; i < LenA; i++ {
		row := make([]float64, colB)
		go populateRows(rowChannel, row, mat1[i], mat2, i, colB)
	}

	for x := 0; x < LenA; x++ {
		select {
		case rowResult := <-rowChannel:
			position := rowResult.index
			finalResult[position] = rowResult.result
		}
	}

	return finalResult
}

func populateRows(c chan RowResult, resultRow []float64, rowA []float64, matrixB [][]float64, position int, colB int) {
	for i, V := range rowA {
		for j := 0; j < colB; j++ {
			valBAtIJ := matrixB[i][j]
			product := V * valBAtIJ
			valAtTempRow := resultRow[j]
			resultRow[j] = product + valAtTempRow
		}
	}
	c <- RowResult{resultRow, position}
}
