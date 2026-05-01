package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeo_CoordinatesByCity_Succes(t *testing.T) {
	result, err := CoordinatesByCity(Moscow)

	require.NoError(t, err)

	assert.Equal(t, 55.751244, result.Latitude)
	assert.Equal(t, 37.618423, result.Longtitude)
}

func TestGeo_CoordinatesByCity_InvalidCityFail(t *testing.T) {
	city := City("invalid-city")
	_, err := CoordinatesByCity(city)

	require.Error(t, err)

	assert.Contains(t, err.Error(), "unknown city")
}
