package confmatrix

import "math"

const LENGTH_DATA_ERROR int = 1

type ConfMatrix struct {
	data        [][]int
	dim         int
	observation int
}

func LoadData(previsional, real []string) (matrix ConfMatrix, err int) {
	if len(previsional) != len(real) {
		err = LENGTH_DATA_ERROR
		return
	}
	var index int
	index = 0

	dataMapping := make(map[string]int)

	for _, value := range previsional {
		if _, ok := dataMapping[value]; !ok {
			dataMapping[value] = index
			index++
		}
	}

	for _, value := range real {
		if _, ok := dataMapping[value]; !ok {
			dataMapping[value] = index
			index++
		}
	}

	matrix = initConfMatrix(len(dataMapping))

	matrix.populate(dataMapping, previsional, real)

	return
}

func initConfMatrix(length int) ConfMatrix {
	var matrix ConfMatrix

	matrix.data = make([][]int, length)
	matrix.dim = length
	matrix.observation = 0

	for i, _ := range matrix.data {
		matrix.data[i] = make([]int, length)
	}

	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			matrix.data[i][j] = 0
		}
	}

	return matrix
}

func (matrix *ConfMatrix) populate(dataMap map[string]int, previsional []string, real []string) {
	for i := 0; i < len(previsional); i++ {
		matrix.data[dataMap[real[i]]][dataMap[previsional[i]]]++
		matrix.observation++
	}
}

func (matrix *ConfMatrix) KCohen() float32 {
	agreementProbability := matrix.agreementProbability()
	errorProbability := matrix.errorProbability()

	return (agreementProbability - errorProbability) / (1 - errorProbability)
}

func (matrix *ConfMatrix) agreementProbability() float32 {
	var sum int = 0
	for i := 0; i < matrix.dim; i++ {
		sum += matrix.data[i][i]
	}
	return float32(sum) / float32(matrix.observation)
}

func (matrix *ConfMatrix) errorProbability() float32 {
	var sum int = 0
	for i := 0; i < matrix.dim; i++ {
		sum += matrix.sumRow(i) * matrix.sumColumn(i)
	}
	return float32(sum) / float32(math.Pow(float64(matrix.observation), 2))
}

func (matrix *ConfMatrix) sumRow(row int) int {
	var sumRow int = 0
	for i := 0; i < matrix.dim; i++ {
		if i == row {
			for j := 0; j < matrix.dim; j++ {
				sumRow += matrix.data[i][j]
			}
		}
	}
	return sumRow
}

func (matrix *ConfMatrix) sumColumn(column int) int {
	var sumColumn int = 0
	for i := 0; i < matrix.dim; i++ {
		if i == column {
			for j := 0; j < matrix.dim; j++ {
				sumColumn += matrix.data[j][i]
			}
		}
	}
	return sumColumn
}
