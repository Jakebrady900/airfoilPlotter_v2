package main

import (
	"airfoilPlotter_v2/M_Matrix"
	"airfoilPlotter_v2/P_Matrix"
	"airfoilPlotter_v2/Parsec"
	"encoding/csv"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"math"
	"os"
	"sync"
)

func main() {
	/*
		-30 < a_TE < 5 [7]
		-0.3 < Y_TE < 0.1 [8]
		0.6 < d2Yl < 0.6 [3]
		The name of the csv and plot will be built by the variables
	*/

	//a_TE := []float64{-30, -22.5, -17.5, -10, -5, 0, 5}
	//Y_TE := []float64{-0.3, -0.25, -0.2, -0.15, -0.1, -0.05, 0, 0.05}
	//d2Yl := []float64{0.0, 0.3, 0.6}
	a_TE := []float64{-30}
	Y_TE := []float64{0.05}
	d2Yl := []float64{0.0}

	wg := new(sync.WaitGroup)
	for i := 0; i < len(a_TE); i++ {
		for j := 0; j < len(Y_TE); j++ {
			for k := 0; k < len(d2Yl); k++ {
				wg.Add(1)
				go generateAirfoil(0.06, 0.34, 0.39, 0.09, -0.287, -0.088, d2Yl[k], 0.000000, Y_TE[j], a_TE[i], 0.00, wg)
			}
		}
	}
	wg.Wait()
	fmt.Println("---DONE---")

}

func generateAirfoil(R_LE, Xu, Xl, Yu, d2Yu, Yl, d2Yl, del_Y_TE, Y_TE, a_TE, b_TE float64, wg *sync.WaitGroup) {
	M := M_Matrix.CreateM(Xu, Xl)
	P := P_Matrix.CreateP(R_LE, Yu, d2Yu, Yl, d2Yl, del_Y_TE, Y_TE, a_TE, b_TE)
	m := M_Matrix.GetInverse(M)
	solution := P_Matrix.Multiply(m, P)
	yUpper, yLower, X := Parsec.GenerateAirfoilUpper(solution)
	isValid := validateArrays(X, yUpper, yLower, Yu, Yl)
	fileName := fmt.Sprintf("%f %f %f ", d2Yl, Y_TE, a_TE)
	plotAirfoil(X, yUpper, yLower, isValid, fileName)
	wg.Done()
}

func plotAirfoil(X, YUpper, YLower []float64, isValid bool, saveName string) {

	ptsUpper := make(plotter.XYs, len(X))
	ptsLower := make(plotter.XYs, len(X))

	for i := range X {
		ptsUpper[i].X = X[i]
		ptsUpper[i].Y = YUpper[i]

		ptsLower[i].X = X[i]
		ptsLower[i].Y = YLower[i]
	}

	p := plot.New()

	p.Title.Text = "Y_Upper and Y_Lower Curves"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	err := plotutil.AddLinePoints(p,
		"Y_Upper", ptsUpper,
		"Y_Lower", ptsLower,
	)
	if err != nil {
		log.Fatalf("plotting failed: %v", err)
	}

	if isValid {
		savePath := "C:\\Users\\jakeb\\GolandProjects\\airfoilPlotter_v2\\Valid\\" + saveName + ".png"
		_ = p.Save(32*vg.Inch, 6*vg.Inch, savePath)
		saveToCSV(X, YUpper, YLower, saveName)
	} else {
		savePath := "C:\\Users\\jakeb\\GolandProjects\\airfoilPlotter_v2\\InValid\\" + saveName + ".png"
		_ = p.Save(32*vg.Inch, 6*vg.Inch, savePath)
	}

}

func saveToCSV(X, YUpper, YLower []float64, saveName string) {
	// Create a csv file
	file, err := os.Create("C:\\Users\\jakeb\\GolandProjects\\airfoilPlotter_v2\\Valid\\" + saveName + ".csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Generate the arrays required to print
	// X repeats twice
	// Y_Upper and Y_Lower are appended to each other
	outputX := make([]float64, (len(X)*2)-1)
	outputY := make([]float64, (len(YUpper) + len(YLower) - 1))
	outputZ := make([]float64, (len(X)*2)-1)

	for i := range X {
		outputX[i] = X[i]
		outputX[i+len(X)-1] = X[i]
		outputY[i] = YUpper[i]
		outputY[i+len(X)-1] = YLower[i]
		outputZ[i] = 0
		outputZ[i+len(X)-1] = 0
	}

	// Write X, Y_Upper, and Y_Lower to the csv file
	for i := range outputX {
		if i == len(outputX)-1 {
			break
		}
		err := writer.Write([]string{fmt.Sprintf("%f", outputX[i]), fmt.Sprintf("%f", outputY[i]), fmt.Sprintf("%f", outputZ[i])})
		if err != nil {
			log.Fatalf("failed writing to file: %s", err)
		}
	}
}

func validateArrays(X, Upper, Lower []float64, Yu, Yl float64) bool {

	if len(X) != len(Upper) || len(X) != len(Lower) {
		fmt.Println("Lengths of X, Y_Upper, and Y_Lower arrays should be the same")
		return false
	}

	//if the max of the first 40% of upper is less than the max of the remaining 60% of upper, then the airfoil is invalid
	break1 := int(float64(len(Upper)) * 0.4)
	section1 := Upper[:break1]
	section2 := Upper[break1:]
	if Max(section1) < Max(section2) {
		fmt.Println("back section higher than front on upper surface")
		return false
	}

	//if the crossover of the upper and lower airfoils is not at the final point, then the airfoil is invalid
	for i := 1; i < len(Upper)-2; i++ {
		if Upper[i] == Lower[i] {
			fmt.Println("crossover not at final point, found at index: ", i)
			return false
		}
	}

	//if the lowest point in the first 40% is less than half of the lowest of the back 60%, then the airfoil is invalid
	section1 = Lower[:break1]
	section2 = Lower[break1:]
	if math.Abs(Min(section1)) < math.Abs(Min(section2))/2 {
		fmt.Println("front section lower than back on lower surface")
		return false
	}

	//if the highest point is higher than 1.5x Yu then the airfoil is invalid
	if Max(Upper) > 1.5*Yu {
		fmt.Println("highest point higher than 1.5x Yu")
		return false
	}

	//if the lowest point is lower than 1.5x Yl then the airfoil is invalid
	if math.Abs(Min(section1)) > math.Abs(Yl*1.5) {
		fmt.Println("lowest point lower than 1.5x Yl")
		return false
	}

	return true
}

func Max(slice []float64) float64 {
	if len(slice) == 0 {
		// Return some default value if the slice is empty
		return math.Inf(-1) // Negative infinity as there's no maximum in an empty slice
	}

	max := slice[0] // Assume the first element is the maximum initially

	for _, value := range slice {
		if value > max {
			max = value
		}
	}

	return max
}

func Min(slice []float64) float64 {
	if len(slice) == 0 {
		// Return some default value if the slice is empty
		return math.Inf(-1) // Negative infinity as there's no maximum in an empty slice
	}

	min := slice[0] // Assume the first element is the maximum initially

	for _, value := range slice {
		if value < min {
			min = value
		}
	}

	return min
}
