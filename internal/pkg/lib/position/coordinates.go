package position

type Coordinates struct {
	Latitude   float64 `json:"latitude" db:"latitude"`
	Longtitude float64 `json:"longtitude" db:"longtitude"`
}

func NewCoordinates(lat float64, long float64) *Coordinates {
	return &Coordinates{
		Latitude:   lat,
		Longtitude: long,
	}
}

func DetermineGeoHashPrecision(distanceKM float32) uint {
	if distanceKM > 625 {
		return 1
	}
	if distanceKM < 625 && distanceKM > 156 {
		return 2
	}
	if distanceKM < 156 && distanceKM > 19.5 {
		return 3
	} else if distanceKM < 19.5 && distanceKM > 4.89 {
		return 4
	} else if distanceKM < 4.89 && distanceKM > 0.61 {
		return 5
	} else if distanceKM < 0.61 && distanceKM > 0.153 {
		return 6
	} else if distanceKM < 0.153 && distanceKM > 0.0191 {
		return 7
	}
	return 8
}
