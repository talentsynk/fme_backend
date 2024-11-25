package utilities

import (
	"strings"
)

const (
	Abia                   string = "abia"
	Adamawa                string = "adamawa"
	AkwaIbom               string = "akwa-ibom"
	Anambra                string = "anambra"
	Bauchi                  string = "bauchi"
	Bayelsa                 string = "bayelsa"
	Benue                  string = "benue"
	Borno                  string = "borno"
	CrossRiver             string = "cross-river"
	Delta                  string = "delta"
	Ebonyi                 string = "ebonyi"
	Edo                    string = "edo"
	Ekiti                  string = "ekiti"
	Enugu                  string = "enugu"
	Gombe                  string = "gombe"
	Imo                    string = "imo"
	Jigawa                 string = "jigawa"
	Kaduna                 string = "kaduna"
	Kano                   string = "kano"
	Katsina                string = "katsina"
	Kebbi                  string = "kebbi"
	Kogi                   string = "kogi"
	Kwara                  string = "kwara"
	Lagos                  string = "lagos"
	Nasarawa               string = "nasarawa"
	Niger                  string = "niger"
	Ogun                   string = "ogun"
	Ondo                   string = "ondo"
	Osun                   string = "osun"
	Oyo                    string = "oyo"
	Plateau                string = "plateau"
	Rivers                 string = "rivers"
	Sokoto                 string = "sokoto"
	Taraba                 string = "taraba"
	Yobe                   string = "yobe"
	Zamfara                string = "zamfara"
	FederalCapitalTerritory string = "federal capital territory"
)

// NigerianStates is a struct that holds all the Nigerian states as constants
type NigerianStates struct{}

func ValidateState(state string) (string, bool) {
	state = strings.ToLower(state)
	switch state {
	case Abia:
	  return Abia, true
	case Adamawa:
	  return Adamawa, true
	case AkwaIbom:
		return AkwaIbom, true
	case Anambra:
		return Anambra, true
	case Bauchi:
		return Bauchi, true
	case Bayelsa:
		return Bayelsa, true
	case Benue:
		return Benue, true
	case Borno:
		return Borno, true
	case CrossRiver:
		return CrossRiver, true
	case Delta:
		return Delta, true
	case Ebonyi:
		return Ebonyi, true
	case Edo:
		return Edo, true
	case Ekiti:
		return Ekiti, true
	case Enugu:
		return Enugu, true
	case Gombe:
		return Gombe, true
	case Imo:
		return Imo, true
	case Jigawa:
		return Jigawa, true
	case Kaduna:
		return Kaduna, true
	case Kano:
		return Kano, true
	case Katsina:
		return Kano, true
	case Kebbi: 
		return Kebbi, true
	case Kogi:
		return Kogi, true
	case Kwara:
		return Kwara, true
	case Lagos:
		return Lagos, true
	case Nasarawa:
		return Nasarawa, true
	case Niger:
		return Niger, true
	case Ogun:
		return Ogun, true
	case Ondo:
		return Ondo, true
	case Osun:
		return Osun, true
	case Oyo:
		return Oyo,true
	case Plateau:
		return Plateau, true
	case Rivers:
		return Rivers, true
	case Sokoto:
		return Sokoto, true
	case Taraba:
		return Taraba, true
	case Yobe:
		return Yobe, true
	case Zamfara:
		return Zamfara, true
	case FederalCapitalTerritory:
		return FederalCapitalTerritory, true
	default:
	  return "",false
	}
  }
