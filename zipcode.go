package zipcode

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Zip struct {
	Code      int
	Latitude  float64
	Longitude float64
	City      string
	State     string
	County    string
	Type      string
	Distance  float64
        LatitudeSin float64
        LatitudeCos float64
}

// ParseCSV parses a string in the form of...
// "zip code","longitude","latitude","city","state","county","type"
func ParseCSV(csv string) (*Zip, error) {
	strs := strings.Split(csv, ",")
	if len(strs) != 7 {
		return nil, errors.New("wrong number of fields")
	}

	zip := new(Zip)

	var err error
	tmp := strings.Replace(strs[0], `"`, "", -1)
	zip.Code, err = strconv.Atoi(tmp)
	if err != nil {
		return nil, err
	}

	tmp = strings.Replace(strs[1], `"`, "", -1)
	if tmp != "" {
		zip.Latitude, err = strconv.ParseFloat(tmp, 64)
		if err != nil {
			fmt.Println(*zip)
			return nil, err
		}
	}

	tmp = strings.Replace(strs[2], `"`, "", -1)
	if tmp != "" {
		zip.Longitude, err = strconv.ParseFloat(tmp, 64)
		if err != nil {
			fmt.Println(*zip)
			return nil, err
		}
	}

	zip.City = strings.Replace(strs[3], `"`, "", -1)
	zip.State = strings.Replace(strs[4], `"`, "", -1)
	zip.County = strings.Replace(strs[5], `"`, "", -1)
	zip.Type = strings.Replace(strs[6], `"`, "", -1)

	return zip, nil
}

// ParseTSV parses a tab-separated-value string in the format
// provided by geonames.org
func ParseTSV(tsv string) (*Zip, error) {
	strs := strings.Split(tsv, "\t")

	if len(strs) != 12 {
		err := fmt.Errorf("found %d fields, expected 12", len(strs))
		return nil, err
	}

	zip := &Zip{}

	var err error
	zip.Code, err = strconv.Atoi(strs[1])
	if err != nil {
		return nil, err
	}

	zip.City = strs[2]
	zip.State = strs[4]
	zip.County = strs[5]

	zip.Latitude, err = strconv.ParseFloat(strs[9], 64)
	if err != nil {
		return nil, err
	}

	zip.Longitude, err = strconv.ParseFloat(strs[10], 64)
	if err != nil {
		return nil, err
	}

	return zip, nil
}

// LoadCSVFile loads a CSV file in the form of...
// "zip code","longitude","latitude","city","state","county","type"
func LoadCSVFile(fileName string) ([]*Zip, error) {
	return loadFile(fileName, ParseCSV)
}

// LoadTSVFile loads a tab-separated-value file in the
// format provided by geonames.org
func LoadTSVFile(fileName string) ([]*Zip, error) {
	return loadFile(fileName, ParseTSV)
}

// Distance calculates the distance between two zip codes.
func Distance(z1, z2 *Zip) float64 {
   if z1.Code == z2.Code {
      return 0.0
   }

   theta := z1.Longitude - z2.Longitude
   d := z1.LatitudeSin * z2.LatitudeSin + z1.LatitudeCos * z2.LatitudeCos * cos(d2r*theta)
   d = math.Acos(d)
   d = r2d * d
   d = d * 60 * 1.1515

   return d
}

// Find takes an integer value for a zip code and returns a pointer to a Zip struct
// or nil if not found.
func Find(zipcode int, zips []*Zip) *Zip {
	for _, zip := range zips {
		if zip.Code == zipcode {
			return zip
		}
	}

	return nil
}

// FindInRadius finds all zip codes within a radius (miles) of zipcode.
func FindInRadius(zipcode int, radius float64, zips []*Zip) []*Zip {
	var found []*Zip
	z1 := Find(zipcode, zips)
	if z1 == nil {
		return found
	}

	// if radius is below threshold distance, skip the
	// slow math and return zipcode
	if radius < 0.1 {
		found = append(found, z1)
		return found
	}

	for _, z2 := range zips {
                d := Distance(z1, z2)

		if d <= radius {
			z2.Distance = d
			found = append(found, z2)
		}
	}

	return found
}

type parser func(s string) (*Zip, error)

func loadFile(fileName string, parse parser) ([]*Zip, error) {
	file, err := os.OpenFile(fileName, 0, 0)
	if err != nil {
		return nil, err
	}

	rdr := bufio.NewReader(file)
	var zips []*Zip

	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			break
		}

		zip, err := parse(line)
		if err != nil {
			return nil, err
		}
                // pre-calculate a couple values for faster distance calculataions
                zip.LatitudeSin = sin(d2r * zip.Latitude)
                zip.LatitudeCos = cos(d2r * zip.Latitude)

		zips = append(zips, zip)
	}
	file.Close()

	return zips, nil
}

const (
	r2d = 180 / math.Pi
	d2r = math.Pi / 180
)

var sin = math.Sin
var cos = math.Cos
