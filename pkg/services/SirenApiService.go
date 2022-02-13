package services

import (
	"github.com/imroc/req/v3"
)

type HttpSirenService struct{}

type DataGouvCompany struct {
	UniteLegale struct {
		Id             int       `json:"id"`
		Siren          string    `json:"siren"`
		// DateCreation   time.Time `json:"date_creation"`
		DateDebut      string    `json:"date_debut"`
		// DateFin        time.Time `json:"date_fin"`
		Denomination   string    `json:"denomination"`
		Etablissements []struct {
			Id                    int       `json:"id"`
			Siren                 string    `json:"siren"`
			Nic                   string    `json:"nic"`
			Siret                 string    `json:"siret"`
			StatutDiffusion       string    `json:"statut_diffusion"`
			// DateCreation          time.Time `json:"date_creation"`
			// DateDernierTraitement time.Time `json:"date_dernier_traitement"`
			CodePostal            string    `json:"code_postal"`
			Longitude             string    `json:"longitude"`
			Latitude              string    `json:"latitude"`
			GeoAdresse            string    `json:"geo_adresse"`
			UniteLegaleId         int       `json:"unite_legale_id"`
		} `json:"etablissements"`
	} `json:"unite_legale"`
}

type Agency struct {
	CodeCompany  string    `json:"siren" csv:"siren"`
	Name         string    `json:"denomination" csv:"denomination"`
	Code         string    `json:"siret" csv:"siret"`
	ZipCode      string    `json:"code_postal" csv:"code_postal"`
	Address      string    `json:"geo_adresse" csv:"geo_adresse"`
}

func (httpSirenService *HttpSirenService) GetAgencies(search SearchSirenApi) ([]Agency, error) {
	zipCode := search.ZipCode
	siren := search.Siren
	var result DataGouvCompany
	resp, err := req.C().R(). // Use R() to create a request
					SetHeader("Accept", "application/json"). // Chainable request settings
					SetPathParam("siren", siren).
					SetResult(&result).
					Get("https://entreprise.data.gouv.fr/api/sirene/v3/unites_legales/{siren}")
	var agencies []Agency
	if resp.IsSuccess() {
		if result.UniteLegale.Siren == siren {
			for _, agency := range result.UniteLegale.Etablissements {
				codePostal := agency.CodePostal
				if codePostal == zipCode {
					agencies = append(agencies, Agency{
						Code:         agency.Siret,
						CodeCompany:  agency.Siren,
						Name:         result.UniteLegale.Denomination,
						Address:      agency.GeoAdresse,
						ZipCode:      codePostal,
					})
				}
			}
		}
	}

	return agencies, err
}
