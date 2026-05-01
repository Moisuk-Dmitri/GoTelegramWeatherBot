package geo

import "fmt"

type Coordinates struct {
	Latitude   float64
	Longtitude float64
}

type City string

const (
	Moscow       City = "moscow"
	StPetersburg City = "stpetersburg"
)

var cityCoordiantes = map[City]Coordinates{
	Moscow:       {Latitude: 55.751244, Longtitude: 37.618423},
	StPetersburg: {Latitude: 59.937500, Longtitude: 30.308611},
}

func CoordinatesByCity(city City) (Coordinates, error) {
	c, ok := cityCoordiantes[city]
	if !ok {
		return Coordinates{}, fmt.Errorf("unknown city: %s", city)
	}

	return c, nil
}
