package main

import (
	"encoding/csv"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Brand struct {
	Name      string //Category
	Count     int    //Count of SKU
	FlavCount int
	FlavName  []string
	Strengths []string //AO Name
}

func main() {
	file1, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file 1: %v\n", err)
		os.Exit(1)
	}
	defer file1.Close()

	file2, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Printf("Error opening file 2: %v\n", err)
		os.Exit(1)
	}
	defer file2.Close()

	reader1 := csv.NewReader(file1)
	reader1.FieldsPerRecord = -1
	reader2 := csv.NewReader(file2)
	reader2.FieldsPerRecord = -1

	rawData1, err := reader1.ReadAll()
	if err != nil {
		fmt.Printf("Error failed to read from file 1: %v\n", err)
		os.Exit(1)
	}
	rawData2, err := reader2.ReadAll()
	if err != nil {
		fmt.Printf("Error failed to read from file 2: %v\n", err)
		os.Exit(1)
	}
	brands := buildBrand(rawData1)

	allBrands := getFlavors(brands, rawData1, rawData2)

	writeExcell(allBrands, os.Args[3])

}

func writeExcell(brands []Brand, outPath string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Totals")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "Name"
	cell = row.AddCell()
	cell.Value = "SKU Count"
	cell = row.AddCell()
	cell.Value = "Flavor Count"

	for v := range brands {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = brands[v].Name
		cell = row.AddCell()
		cell.Value = strconv.Itoa(brands[v].Count)
		cell = row.AddCell()
		cell.Value = strconv.Itoa(brands[v].FlavCount)

		if v == len(brands)-1 {
			row = sheet.AddRow()
			cell = row.AddCell()
			cell = row.AddCell()
			cell.Value = "Total SKU"
			cell = row.AddCell()
			cell.Value = "Total Flavors"
			row = sheet.AddRow()
			cell = row.AddCell()

			count := 0
			flvCount := 0
			for i := range brands {
				count += brands[i].Count
				flvCount += brands[i].FlavCount
			}
			cell = row.AddCell()
			cell.Value = strconv.Itoa(count)
			cell = row.AddCell()
			cell.Value = strconv.Itoa(flvCount)
		}
	}

	sheet, err = file.AddSheet("Strengths")
	for i := range brands {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = brands[i].Name
		for j := range brands[i].Strengths {
			cell = row.AddCell()
			cell.Value = brands[i].Strengths[j]
		}
	}

	err = file.Save(outPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getFlavors(brands []string, rawData1, rawData2 [][]string) []Brand {
	allBrands := []Brand{}
	for i := 0; i < len(brands); i++ {
		flavName := []string{}
		strengths := []string{}
		var count int
		for _, data := range rawData1 {
			if data[3] == brands[i] {
				if data[48] == "0" {
					flavName = append(flavName, data[2])
					flavStrength, c := getStrengths(data[0], rawData2, strengths)
					count += c
					strengths = parseDupStrength(flavStrength)
				}
			}
		}
		allBrands = append(allBrands, Brand{
			Name:      brands[i],
			Count:     len(flavName) * len(strengths),
			FlavCount: len(flavName),
			FlavName:  flavName,
			Strengths: strengths,
		})
	}
	return allBrands
}

func parseDupStrength(strengths []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range strengths {
		spaceless := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, strengths[v])
		if encountered[spaceless] != true {
			encountered[spaceless] = true
			result = append(result, strengths[v])
		}
	}
	return result
}

func getStrengths(flavID string, rawData2 [][]string, strength []string) ([]string, int) {
	count := 0
	for _, data := range rawData2 {
		if flavID == data[0] {
			strength = append(strength, data[4])
			count++
		}
	}
	return strength, count
}

func buildBrand(rawData1 [][]string) []string {
	brands := []string{}

	brands = append(brands, rawData1[1][3])
	for _, data1 := range rawData1 {
	InnerLoop:
		for i := 0; i < len(brands); i++ {
			if data1[3] == brands[i] {
				break InnerLoop
			}
			if i == len(brands)-1 {
				brands = append(brands, data1[3])
			}
		}
	}
	return brands
}
