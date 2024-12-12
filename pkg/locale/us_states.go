package goliblocale

import "strings"

type USStateCodeLen2 string

const (
	USStateCodeLen2Alabama                = USStateCodeLen2("AL")
	USStateCodeLen2Alaska                 = USStateCodeLen2("AK")
	USStateCodeLen2Arizona                = USStateCodeLen2("AZ")
	USStateCodeLen2Arkansas               = USStateCodeLen2("AR")
	USStateCodeLen2AmericanSamoa          = USStateCodeLen2("AS")
	USStateCodeLen2California             = USStateCodeLen2("CA")
	USStateCodeLen2Colorado               = USStateCodeLen2("CO")
	USStateCodeLen2Connecticut            = USStateCodeLen2("CT")
	USStateCodeLen2Delaware               = USStateCodeLen2("DE")
	USStateCodeLen2DistrictOfColumbia     = USStateCodeLen2("DC")
	USStateCodeLen2Florida                = USStateCodeLen2("FL")
	USStateCodeLen2Georgia                = USStateCodeLen2("GA")
	USStateCodeLen2Guam                   = USStateCodeLen2("GU")
	USStateCodeLen2Hawaii                 = USStateCodeLen2("HI")
	USStateCodeLen2Idaho                  = USStateCodeLen2("ID")
	USStateCodeLen2Illinois               = USStateCodeLen2("IL")
	USStateCodeLen2Indiana                = USStateCodeLen2("IN")
	USStateCodeLen2Iowa                   = USStateCodeLen2("IA")
	USStateCodeLen2Kansas                 = USStateCodeLen2("KS")
	USStateCodeLen2Kentucky               = USStateCodeLen2("KY")
	USStateCodeLen2Louisiana              = USStateCodeLen2("LA")
	USStateCodeLen2Maine                  = USStateCodeLen2("ME")
	USStateCodeLen2Maryland               = USStateCodeLen2("MD")
	USStateCodeLen2Massachusetts          = USStateCodeLen2("MA")
	USStateCodeLen2Michigan               = USStateCodeLen2("MI")
	USStateCodeLen2Minnesota              = USStateCodeLen2("MN")
	USStateCodeLen2Mississippi            = USStateCodeLen2("MS")
	USStateCodeLen2Missouri               = USStateCodeLen2("MO")
	USStateCodeLen2Montana                = USStateCodeLen2("MT")
	USStateCodeLen2Nebraska               = USStateCodeLen2("NE")
	USStateCodeLen2Nevada                 = USStateCodeLen2("NV")
	USStateCodeLen2NewHampshire           = USStateCodeLen2("NH")
	USStateCodeLen2NewJersey              = USStateCodeLen2("NJ")
	USStateCodeLen2NewMexico              = USStateCodeLen2("NM")
	USStateCodeLen2NewYork                = USStateCodeLen2("NY")
	USStateCodeLen2NorthCarolina          = USStateCodeLen2("NC")
	USStateCodeLen2NorthDakota            = USStateCodeLen2("ND")
	USStateCodeLen2NorthernMarianaIslands = USStateCodeLen2("MP") //doubt
	USStateCodeLen2Ohio                   = USStateCodeLen2("OH")
	USStateCodeLen2Oklahoma               = USStateCodeLen2("OK")
	USStateCodeLen2Oregon                 = USStateCodeLen2("OR")
	USStateCodeLen2Pennsylvania           = USStateCodeLen2("PA")
	USStateCodeLen2PuertoRico             = USStateCodeLen2("PR")
	USStateCodeLen2RhodeIsland            = USStateCodeLen2("RI")
	USStateCodeLen2SouthCarolina          = USStateCodeLen2("SC")
	USStateCodeLen2SouthDakota            = USStateCodeLen2("SD")
	USStateCodeLen2Tennessee              = USStateCodeLen2("TN")
	USStateCodeLen2Texas                  = USStateCodeLen2("TX")
	USStateCodeLen2TrustTerritories       = USStateCodeLen2("TT")
	USStateCodeLen2Utah                   = USStateCodeLen2("UT")
	USStateCodeLen2Vermont                = USStateCodeLen2("VT")
	USStateCodeLen2Virginia               = USStateCodeLen2("VA")
	USStateCodeLen2VirginIslands          = USStateCodeLen2("VI")
	USStateCodeLen2Washington             = USStateCodeLen2("WA")
	USStateCodeLen2WestVirginia           = USStateCodeLen2("WV")
	USStateCodeLen2Wisconsin              = USStateCodeLen2("WI")
	USStateCodeLen2Wyoming                = USStateCodeLen2("WY")
)

type USStateName string

const (
	USStateNameAlabama                = USStateName("Alabama")
	USStateNameAlaska                 = USStateName("Alaska")
	USStateNameArizona                = USStateName("Arizona")
	USStateNameArkansas               = USStateName("Arkansas")
	USStateNameAmericanSamoa          = USStateName("American Samoa")
	USStateNameCalifornia             = USStateName("California")
	USStateNameColorado               = USStateName("Colorado")
	USStateNameConnecticut            = USStateName("Connecticut")
	USStateNameDelaware               = USStateName("Delaware")
	USStateNameDistrictOfColumbia     = USStateName("District Of Columbia")
	USStateNameFlorida                = USStateName("Florida")
	USStateNameGeorgia                = USStateName("Georgia")
	USStateNameGuam                   = USStateName("Guam")
	USStateNameHawaii                 = USStateName("Hawaii")
	USStateNameIdaho                  = USStateName("Idaho")
	USStateNameIllinois               = USStateName("Illinois")
	USStateNameIndiana                = USStateName("Indiana")
	USStateNameIowa                   = USStateName("Iowa")
	USStateNameKansas                 = USStateName("Kansas")
	USStateNameKentucky               = USStateName("Kentucky")
	USStateNameLouisiana              = USStateName("Louisiana")
	USStateNameMaine                  = USStateName("Maine")
	USStateNameMaryland               = USStateName("Maryland")
	USStateNameMassachusetts          = USStateName("Massachusetts")
	USStateNameMichigan               = USStateName("Michigan")
	USStateNameMinnesota              = USStateName("Minnesota")
	USStateNameMississippi            = USStateName("Mississippi")
	USStateNameMissouri               = USStateName("Missouri")
	USStateNameMontana                = USStateName("Montana")
	USStateNameNebraska               = USStateName("Nebraska")
	USStateNameNevada                 = USStateName("Nevada")
	USStateNameNewHampshire           = USStateName("New Hampshire")
	USStateNameNewJersey              = USStateName("New Jersey")
	USStateNameNewMexico              = USStateName("New Mexico")
	USStateNameNewYork                = USStateName("New York")
	USStateNameNorthCarolina          = USStateName("North Carolina")
	USStateNameNorthDakota            = USStateName("North Dakota")
	USStateNameNorthernMarianaIslands = USStateName("Northern Mariana Islands")
	USStateNameOhio                   = USStateName("Ohio")
	USStateNameOklahoma               = USStateName("Oklahoma")
	USStateNameOregon                 = USStateName("Oregon")
	USStateNamePennsylvania           = USStateName("Pennsylvania")
	USStateNamePuertoRico             = USStateName("Puerto Rico")
	USStateNameRhodeIsland            = USStateName("Rhode Island")
	USStateNameSouthCarolina          = USStateName("SouthCarolina")
	USStateNameSouthDakota            = USStateName("SouthDakota")
	USStateNameTennessee              = USStateName("Tennessee")
	USStateNameTexas                  = USStateName("Texas")
	USStateNameTrustTerritories       = USStateName("Trust Territories")
	USStateNameUtah                   = USStateName("Utah")
	USStateNameVermont                = USStateName("Vermont")
	USStateNameVirginia               = USStateName("Virginia")
	USStateNameVirginIslands          = USStateName("Virgin Islands")
	USStateNameWashington             = USStateName("Washington")
	USStateNameWestVirginia           = USStateName("West Virginia")
	USStateNameWisconsin              = USStateName("Wisconsin")
	USStateNameWyoming                = USStateName("Wyoming")
)

var USStateNameToCode2Len = map[USStateName]USStateCodeLen2{
	USStateNameAlabama:                USStateCodeLen2Alabama,
	USStateNameAlaska:                 USStateCodeLen2Alaska,
	USStateNameArizona:                USStateCodeLen2Arizona,
	USStateNameArkansas:               USStateCodeLen2Arkansas,
	USStateNameAmericanSamoa:          USStateCodeLen2AmericanSamoa,
	USStateNameCalifornia:             USStateCodeLen2California,
	USStateNameColorado:               USStateCodeLen2Colorado,
	USStateNameConnecticut:            USStateCodeLen2Connecticut,
	USStateNameDelaware:               USStateCodeLen2Delaware,
	USStateNameDistrictOfColumbia:     USStateCodeLen2DistrictOfColumbia,
	USStateNameFlorida:                USStateCodeLen2Florida,
	USStateNameGeorgia:                USStateCodeLen2Georgia,
	USStateNameGuam:                   USStateCodeLen2Guam,
	USStateNameHawaii:                 USStateCodeLen2Hawaii,
	USStateNameIdaho:                  USStateCodeLen2Idaho,
	USStateNameIllinois:               USStateCodeLen2Illinois,
	USStateNameIndiana:                USStateCodeLen2Indiana,
	USStateNameIowa:                   USStateCodeLen2Iowa,
	USStateNameKansas:                 USStateCodeLen2Kansas,
	USStateNameKentucky:               USStateCodeLen2Kentucky,
	USStateNameLouisiana:              USStateCodeLen2Louisiana,
	USStateNameMaine:                  USStateCodeLen2Maine,
	USStateNameMaryland:               USStateCodeLen2Maryland,
	USStateNameMassachusetts:          USStateCodeLen2Massachusetts,
	USStateNameMichigan:               USStateCodeLen2Michigan,
	USStateNameMinnesota:              USStateCodeLen2Minnesota,
	USStateNameMississippi:            USStateCodeLen2Mississippi,
	USStateNameMissouri:               USStateCodeLen2Missouri,
	USStateNameMontana:                USStateCodeLen2Montana,
	USStateNameNebraska:               USStateCodeLen2Nebraska,
	USStateNameNevada:                 USStateCodeLen2Nevada,
	USStateNameNewHampshire:           USStateCodeLen2NewHampshire,
	USStateNameNewJersey:              USStateCodeLen2NewJersey,
	USStateNameNewMexico:              USStateCodeLen2NewMexico,
	USStateNameNewYork:                USStateCodeLen2NewYork,
	USStateNameNorthCarolina:          USStateCodeLen2NorthCarolina,
	USStateNameNorthDakota:            USStateCodeLen2NorthDakota,
	USStateNameNorthernMarianaIslands: USStateCodeLen2NorthernMarianaIslands,
	USStateNameOhio:                   USStateCodeLen2Ohio,
	USStateNameOklahoma:               USStateCodeLen2Oklahoma,
	USStateNameOregon:                 USStateCodeLen2Oregon,
	USStateNamePennsylvania:           USStateCodeLen2Pennsylvania,
	USStateNamePuertoRico:             USStateCodeLen2PuertoRico,
	USStateNameRhodeIsland:            USStateCodeLen2RhodeIsland,
	USStateNameSouthCarolina:          USStateCodeLen2SouthCarolina,
	USStateNameSouthDakota:            USStateCodeLen2SouthDakota,
	USStateNameTennessee:              USStateCodeLen2Tennessee,
	USStateNameTexas:                  USStateCodeLen2Texas,
	USStateNameTrustTerritories:       USStateCodeLen2TrustTerritories,
	USStateNameUtah:                   USStateCodeLen2Utah,
	USStateNameVermont:                USStateCodeLen2Vermont,
	USStateNameVirginia:               USStateCodeLen2Virginia,
	USStateNameVirginIslands:          USStateCodeLen2VirginIslands,
	USStateNameWashington:             USStateCodeLen2Washington,
	USStateNameWestVirginia:           USStateCodeLen2WestVirginia,
	USStateNameWisconsin:              USStateCodeLen2Wisconsin,
	USStateNameWyoming:                USStateCodeLen2Wyoming,
}

func GetUSStateCodeLen2(state string) string {

	for stateName, stateCode := range USStateNameToCode2Len {
		usStateParts := strings.Fields(string(stateName))
		stateParts := strings.Fields(state)

		if len(usStateParts) != len(stateParts) {
			continue
		}
		found := true
		for i := 0; i < len(stateParts); i++ {
			if !strings.EqualFold(usStateParts[i], stateParts[i]) {
				found = false
				break
			}
		}

		if found {
			return string(stateCode)
		}

	}

	return state
}
