package utils

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
)

func LoadCsv(csvfile string) ([]float64, []float64, error) {
	var vals0 []float64
	var vals1 []float64
	var reader *csv.Reader
	var row []string

	read_file, err := os.Open(csvfile)
	if err != nil {
		return vals0, vals1, err
	}
	defer read_file.Close()

	reader = csv.NewReader(read_file)
	// Skip a line.
	_, err = reader.Read()
	if err != nil {
		return vals0, vals1, err
	}

	for {
		row, err = reader.Read()
		if err == io.EOF {
			return vals0, vals1, nil
		}
		val0, _ := strconv.ParseFloat(strings.TrimSpace(row[0]), 64)
		vals0 = append(vals0, val0)
		val1, _ := strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
		vals1 = append(vals1, val1)
	}
}

func LoadDesignSpectrum(csvfile string) ([]float64, []float64, error) {
	var t []float64
	var sa []float64
	var err error

	t, sa, err = LoadCsv(csvfile)
	if err != nil {
		return t, sa, err
	}
	return t, sa, nil
}
