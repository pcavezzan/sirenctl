package services

import (
	"encoding/csv"
	"errors"
	"github.com/gocarina/gocsv"
	"io"
	"os"
)

type SearchSirenApi struct {
	Siren string
	ZipCode string
}

type SearchSirenApiParser interface {
	Parse() ([]SearchSirenApi, error)
}

type argumentZipCodeParser struct {
	siren, zipCode string
}

func NewArgumentZipCodeParser(siren string, zipCode string) SearchSirenApiParser {
	return &argumentZipCodeParser{siren: siren, zipCode: zipCode}
}

func (a *argumentZipCodeParser) Parse() ([]SearchSirenApi, error) {
	if a.siren == "" {
		return nil, errors.New("siren is required")
	}

	if a.zipCode == "" {
		return nil, errors.New("zipCode is required")
	}

	return []SearchSirenApi{
		{
			Siren:   a.siren,
			ZipCode: a.zipCode,
		},
	}, nil
}

type csvFileSirenZipCodeParser struct {
	filePath string
}

type inputCsvContent struct {
	Siren string `csv:"siren"`
	ZipCode string `csv:"code_postal"`
}

func NewCsvFileSirenZipCodeParser(filePath string, separator rune) SearchSirenApiParser {
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = separator
		return r // Allows use pipe as delimiter
	})

	return &csvFileSirenZipCodeParser{filePath: filePath}
}

func (c *csvFileSirenZipCodeParser) Parse() ([]SearchSirenApi, error) {
	if c.filePath == "" {
		return nil, errors.New("file path is required")
	}

	sirenZipCodeFile, err := os.OpenFile(c.filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer sirenZipCodeFile.Close()

	var inputCsvContent []*inputCsvContent
	if err = gocsv.UnmarshalFile(sirenZipCodeFile, &inputCsvContent); err != nil {
		return nil, err
	}

	searchSirenApi := make([]SearchSirenApi, len(inputCsvContent))
	for _, content := range inputCsvContent {
		searchSirenApi = append(searchSirenApi, SearchSirenApi{
			Siren:   content.Siren,
			ZipCode: content.ZipCode,
		})
	}
	return searchSirenApi, nil
}
