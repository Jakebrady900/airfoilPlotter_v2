package M_Matrix

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
)

func CreateM(Xu, Xl float64) [][]float64 {

	M := doCreate()
	doRow1(M)
	doRow2(M, Xu) //Xu
	doRow3(M, Xu) //Xu
	doRow4(M, Xl) //Xl
	doRow5(M, Xl) //Xl
	doRow6(M)
	doRow7(M)
	doRow8(M)
	doRow9(M)
	doRow10(M, Xu)
	doRow11(M, Xl)
	doRow12(M)

	return M
}

func doCreate() [][]float64 {
	M := make([][]float64, 12)
	for x := 0; x < 12; x++ {
		M[x] = []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}
	return M
}

func doRow1(M [][]float64) {
	M[0][0] = 1
}

func doRow2(M [][]float64, Xu float64) {
	pows := []float64{0.5, 1.5, 2.5, 3.5, 4.5, 5.5}
	for n := 0; n < 6; n++ {
		M[1][n] = math.Pow(Xu, pows[n])
	}
}
func doRow3(M [][]float64, Xu float64) {
	coeffs := []float64{-0.25, 0.75, 3.75, 8.75, 15.75, 24.75}
	pows := []float64{-1.5, -0.5, 0.5, 1.5, 2.5, 3.5}
	for n := 0; n < 6; n++ {
		M[2][n] = coeffs[n] * (math.Pow(Xu, pows[n]))
	}
}
func doRow4(M [][]float64, Xl float64) {
	pows := []float64{0.5, 1.5, 2.5, 3.5, 4.5, 5.5}
	for n := 0; n < 6; n++ {
		M[3][n+6] = math.Pow(Xl, pows[n])
	}
}
func doRow5(M [][]float64, Xl float64) {
	coeffs := []float64{-0.25, 0.75, 3.75, 8.75, 15.75, 24.75}
	pows := []float64{-1.5, -0.5, 0.5, 1.5, 2.5, 3.5}
	for n := 0; n < 6; n++ {
		M[4][n+6] = coeffs[n] * (math.Pow(Xl, pows[n]))
	}
}
func doRow6(M [][]float64) {
	for n := 0; n < 6; n++ {
		M[5][n] = 1
	}
	for n := 6; n < 12; n++ {
		M[5][n] = -1
	}
}
func doRow7(M [][]float64) {
	for n := 0; n < 12; n++ {
		M[6][n] = 0.5
	}
}
func doRow8(M [][]float64) {
	nums := []float64{1, 3, 5, 7, 9, 11}
	for n := 0; n < 6; n++ {
		M[7][n] = nums[n] / 2
	}
}
func doRow9(M [][]float64) {
	nums := []float64{1, 3, 5, 7, 9, 11}
	for n := 0; n < 6; n++ {
		M[8][n+6] = nums[n] / 2
	}
}
func doRow10(M [][]float64, Xu float64) {
	coeffs := []float64{0.5, 1.5, 2.5, 3.5, 4.5, 5.5}
	pows := []float64{-0.5, 0.5, 1.5, 2.5, 3.5, 4.5}
	for n := 0; n < 6; n++ {
		M[9][n] = coeffs[n] * (math.Pow(Xu, pows[n]))
	}
}
func doRow11(M [][]float64, Xl float64) {
	coeffs := []float64{0.5, 1.5, 2.5, 3.5, 4.5, 5.5}
	pows := []float64{-0.5, 0.5, 1.5, 2.5, 3.5, 4.5}
	for n := 0; n < 6; n++ {
		M[10][n+6] = coeffs[n] * (math.Pow(Xl, pows[n]))
	}
}
func doRow12(M [][]float64) {
	M[11][6] = 1
}

func PrintMatrix(M [][]float64) {
	for _, row := range M {
		for _, val := range row {
			fmt.Printf("%6.3f\t  ", val)
		}
		fmt.Println()
	}
}

func GetInverse(M [][]float64) [][]float64 {
	// Convert [][]float64 to mat.Dense
	rows, cols := len(M), len(M[0])
	matrix := mat.NewDense(rows, cols, nil)

	// Fill the mat.Dense matrix with values from the slice
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			matrix.Set(i, j, M[i][j])
		}
	}

	// Calculate the inverse matrix
	var inv mat.Dense
	err := inv.Inverse(matrix)
	if err != nil {
		log.Fatal(err)
	}

	// Get the inverse as [][]float64
	var invSlice [][]float64
	invRows, invCols := inv.Dims()
	invSlice = make([][]float64, invRows)
	for i := 0; i < invRows; i++ {
		invSlice[i] = make([]float64, invCols)
		for j := 0; j < invCols; j++ {
			invSlice[i][j] = inv.At(i, j)
		}
	}

	return invSlice
}
