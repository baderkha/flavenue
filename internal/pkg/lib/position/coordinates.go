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
	if distanceKM > 4992.6 {
		return 1
	} else if distanceKM < 4992.6 && distanceKM > 624.1 {
		return 2
	} else if distanceKM < 624.1 && distanceKM > 156 {
		return 3
	} else if distanceKM < 156 && distanceKM > 19.5 {
		return 4
	} else if distanceKM < 19.5 && distanceKM > 4.9 {
		return 5
	} else if distanceKM < 4.9 && distanceKM > 0.6094 {
		return 6
	} else if distanceKM < 0.6094 && distanceKM > 0.1529 {
		return 7
	}
	return 8
}
