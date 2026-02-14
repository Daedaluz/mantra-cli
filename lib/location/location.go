package location

import (
	"math"
	"strings"

	"github.com/daedaluz/mantra-cli/lib/grpc/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	base32 = "0123456789bcdefghjkmnpqrstuvwxyz"
)

// ValidateLocation validates that a location has at least one location identifier
// and that geohash and coordinates are mutually exclusive
func ValidateLocation(loc *common.Location) error {
	if loc == nil {
		return nil
	}

	hasLocation := false
	if loc.Ip != "" {
		hasLocation = true
	}
	if loc.Geohash != "" {
		hasLocation = true
	}
	if loc.Position != nil {
		hasLocation = true
	}

	if !hasLocation {
		return status.Errorf(codes.InvalidArgument, "location must have at least one of: ip, geohash, or position")
	}

	// Check mutual exclusivity between geohash and coordinates
	if loc.Geohash != "" && loc.Position != nil {
		return status.Errorf(codes.InvalidArgument, "location cannot have both geohash and coordinates - they are mutually exclusive")
	}

	return nil
}

// GetLocationCoordinates extracts coordinates from a location, converting from geohash if necessary.
// For geohashes, accuracy is set to the distance from the center to the corner of the geohash cell.
func GetLocationCoordinates(loc *common.Location) (*common.LongLatitude, error) {
	if loc == nil {
		return nil, nil
	}

	// If coordinates are already provided, return them
	if loc.Position != nil {
		return loc.Position, nil
	}

	// If geohash is provided, decode it to coordinates with accuracy
	if loc.Geohash != "" {
		latMin, latMax, lonMin, lonMax, err := DecodeGeohashBounds(loc.Geohash)
		if err != nil {
			return nil, err
		}
		lat := (latMin + latMax) / 2
		lon := (lonMin + lonMax) / 2
		// Accuracy is the distance from center to corner of the bounding box
		accuracy := CalculateDistance(lat, lon, latMax, lonMax)
		return &common.LongLatitude{
			Latitude:  lat,
			Longitude: lon,
			Accuracy:  accuracy,
		}, nil
	}

	// No location coordinates available
	return nil, nil
}

// CalculateDistance calculates the distance between two points using the Haversine formula
// Returns distance in meters
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth's radius in meters

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// ValidateGeohash validates a geohash string
func ValidateGeohash(geohash string) error {
	if geohash == "" {
		return nil
	}

	// Basic geohash validation - should only contain base32 characters
	validChars := "0123456789bcdefghjkmnpqrstuvwxyz"
	for _, char := range strings.ToLower(geohash) {
		if !strings.ContainsRune(validChars, char) {
			return status.Errorf(codes.InvalidArgument, "invalid geohash character: %c", char)
		}
	}

	return nil
}

// ValidateCoordinates validates latitude and longitude values
func ValidateCoordinates(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return status.Errorf(codes.InvalidArgument, "latitude must be between -90 and 90 degrees")
	}
	if lon < -180 || lon > 180 {
		return status.Errorf(codes.InvalidArgument, "longitude must be between -180 and 180 degrees")
	}
	return nil
}

// DecodeGeohashBounds decodes a geohash string to its bounding box
func DecodeGeohashBounds(geohash string) (latMin, latMax, lonMin, lonMax float64, err error) {
	if geohash == "" {
		return 0, 0, 0, 0, status.Errorf(codes.InvalidArgument, "geohash cannot be empty")
	}

	geohash = strings.ToLower(geohash)

	// Validate geohash characters
	for _, char := range geohash {
		if !strings.ContainsRune(base32, char) {
			return 0, 0, 0, 0, status.Errorf(codes.InvalidArgument, "invalid geohash character: %c", char)
		}
	}

	latMin, latMax = -90.0, 90.0
	lonMin, lonMax = -180.0, 180.0

	even := true
	for _, char := range geohash {
		idx := strings.IndexRune(base32, char)
		if idx == -1 {
			return 0, 0, 0, 0, status.Errorf(codes.InvalidArgument, "invalid geohash character: %c", char)
		}

		for i := 4; i >= 0; i-- {
			bit := (idx >> i) & 1
			if even {
				// longitude
				if bit == 1 {
					lonMin = (lonMin + lonMax) / 2
				} else {
					lonMax = (lonMin + lonMax) / 2
				}
			} else {
				// latitude
				if bit == 1 {
					latMin = (latMin + latMax) / 2
				} else {
					latMax = (latMin + latMax) / 2
				}
			}
			even = !even
		}
	}

	return latMin, latMax, lonMin, lonMax, nil
}
