package services

import (
	"encoding/csv"
	"github.com/gocarina/gocsv"
	"io"
)

type AgencyFormatter interface {
	Format(agencies []Agency) (string, error)
}

type csvAgencyFormatter struct{}

func NewCsvAgencyFormatter(separator rune) AgencyFormatter {
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = separator
		return gocsv.NewSafeCSVWriter(writer)
	})
	return &csvAgencyFormatter{}
}

type csvContent struct {
	Siren   string `csv:"siren"`
	ZipCode string `csv:"code_postal"`
	Siret    string `csv:"siret"`
}

func (c *csvAgencyFormatter) Format(agencies []Agency) (string, error) {
	var rows []csvContent
	for _, agency := range agencies {
		csvRow := csvContent{
			Siren:   agency.CodeCompany,
			Siret:   agency.Code,
			ZipCode: agency.ZipCode,
		}
		rows = append(rows, csvRow)
	}

	csvContent, err := gocsv.MarshalString(&rows)
	if err != nil {
		return "", err
	}
	return csvContent, nil
}
