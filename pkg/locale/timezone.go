package goliblocale

import (
	"errors"
	"time"
)

type TimeZone string

const (
	TimeZoneUnknown            = TimeZone("")
	TimeZoneAmericaHalifax     = TimeZone("America/Halifax")
	TimeZoneAmericaNewYork     = TimeZone("America/New_York")
	TimeZoneAmericaChicago     = TimeZone("America/Chicago")
	TimeZoneAmericaDenver      = TimeZone("America/Denver")
	TimeZoneAmericaPhoenix     = TimeZone("America/Phoenix")
	TimeZoneAmericaLosAngeles  = TimeZone("America/Los_Angeles")
	TimeZoneAmericaAnchorage   = TimeZone("America/Anchorage")
	TimeZoneAmericaAdak        = TimeZone("America/Adak")
	TimeZonePacificHonolulu    = TimeZone("Pacific/Honolulu")
	TimeZonePacificPagoPago    = TimeZone("Pacific/Pago_Pago")
	TimeZoneAsiaTokyo          = TimeZone("Asia/Tokyo")
	TimeZonePacificGuam        = TimeZone("Pacific/Guam")
	TimeZonePacificGuadalcanal = TimeZone("Pacific/Guadalcanal")
	TimeZonePacificWake        = TimeZone("Pacific/Wake")
)

func (tz TimeZone) Location() (*time.Location, error) {
	if tz == TimeZoneUnknown {
		return nil, errors.New("invalid time zone")
	}

	return time.LoadLocation(string(tz))
}
