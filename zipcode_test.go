package zipcode

import "testing"

func TestParseCSV(t *testing.T) {
	z, err := ParseCSV(`"00501","+40.922326","-072.637078","HOLTSVILLE","NY","SUFFOLK","UNIQUE"`)
	if err != nil {
		t.Error(err)
	} else if z.Code != 501 {
		t.Errorf("z.Code = %v want 501", z.Code)
	} else if !floatEqual(z.Latitude, 40.922326, 0.0000009) {
		t.Errorf("z.Latitude = %v want 40.922326", z.Latitude)
	} else if !floatEqual(z.Longitude, -72.637078, 0.0000009) {
		t.Errorf("z.Longitude = %v want -72.637078", z.Longitude)
	} else if z.City != "HOLTSVILLE" {
		t.Errorf("z.State = %v want HOLTSVILLE", z.City)
	} else if z.State != "NY" {
		t.Errorf("z.State = %v want NY", z.State)
	} else if z.County != "SUFFOLK" {
		t.Errorf("z.County = %v want SUFFOLK", z.County)
	} else if z.Type != "UNIQUE" {
		t.Errorf("z.Type = %v want UNIQUE", z.Type)
	}
}

func TestParseTSV(t *testing.T) {
	z, err := ParseTSV("US\t00501\tHoltsville\tNew York\tNY\tSuffolk\t103\t\t\t40.9223\t-72.6371\t")
	if err != nil {
		t.Error(err)
	} else if z.Code != 501 {
		t.Errorf("z.Code = %v want 501", z.Code)
	} else if !floatEqual(z.Latitude, 40.9223, 0.0000009) {
		t.Errorf("z.Latitude = %v want 40.9223", z.Latitude)
	} else if !floatEqual(z.Longitude, -72.6371, 0.0000009) {
		t.Errorf("z.Longitude = %v want -72.6371", z.Longitude)
	} else if z.City != "Holtsville" {
		t.Errorf("z.State = %v want Holtsville", z.City)
	} else if z.State != "NY" {
		t.Errorf("z.State = %v want NY", z.State)
	} else if z.County != "Suffolk" {
		t.Errorf("z.County = %v want Suffolk", z.County)
	}
}

func TestLoadTSVFile(t *testing.T) {
	zips, err := LoadTSVFile("US.txt")
	if err != nil {
		t.Error(err)
		return
	}

	if len(zips) != 43628 {
		t.Errorf("len(zips) = %d want 42741", len(zips))
	}
}

func TestDistance(t *testing.T) {
	delta := 0.01
	d := Distance(0, 0, 0, 0)
	if !floatEqual(d, 0.0, delta) {
		t.Errorf("Distance(0, 0, 0, 0) = %f want 0 +/- %f", d, delta)
	}

	d = Distance(40.922326, -72.637078, 40.922326, -72.637078)
	if !floatEqual(d, 0.0, delta) {
		t.Errorf("Distance(40.922326, -72.637078, 40.922326, -72.637078) = %f want 0 +/- %f", d, delta)
	}

	d = Distance(40.922326, -72.637078, 35.688136, -80.819825)
	if !floatEqual(d, 571.90, delta) {
		t.Errorf("Distance(40.922326, -72.637078, 40.922326, -72.637078) = %f want 571.90 +/- %f", d, delta)
	}
}

func TestFind(t *testing.T) {
	zips, err := LoadCSVFile("zip_codes.csv")
	if err != nil {
		t.Error(err)
		return
	}

	zip := Find(28115, zips)
	if zip == nil {
		t.Errorf("Find(28115, zips) = nil want a valid *Zip")
		return
	}

	if zip.Code != 28115 {
		t.Errorf("Find(28115, zips) found the wrong zip! zip.Code = %d want 28115", zip.Code)
	}
}

func TestFindInRadius(t *testing.T) {
	zips, err := LoadCSVFile("zip_codes.csv")
	if err != nil {
		t.Error(err)
		return
	}

	found := FindInRadius(28115, 10, zips)
	want := 8
	if len(found) != want {
		t.Errorf("FindInRadius(28115, 10, zips) found %d zip codes.  Want %d", len(found), want)
	}

	zips, err = LoadTSVFile("US.txt")
	if err != nil {
		t.Error(err)
		return
	}

	found = FindInRadius(28115, 10, zips)
	want = 13
	if len(found) != want {
		t.Errorf("FindInRadius(28115, 10, zips) found %d zip codes.  Want %d", len(found), want)
	}
}

func floatEqual(a, b, delta float64) bool {
	if (b-delta) <= a && a <= (b+delta) {
		return true
	}
	return false
}
