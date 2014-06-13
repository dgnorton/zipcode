package zipcode

import "testing"

func TestParseTSV(t *testing.T) {
	z, err := ParseTSV("US\t00501\tHoltsville\tNew York\tNY\tSuffolk\t103\t\t\t40.9223\t-72.6371\t")
	if err != nil {
		t.Error(err)
	} else if z.Code != "00501" {
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

	err = nil
	z, err = ParseTSV("BOGUS\tUS\t00501\tHoltsville\tNew York\tNY\tSuffolk\t103\t\t\t40.9223\t-72.6371\t")
	if err == nil {
		t.Errorf("expected an error due to too many fields")
	}

	err = nil
	z, err = ParseTSV("US\t00501\tHoltsville\tNew York\tNY\tSuffolk\t103\t\t\tBadLatitude\t-72.6371\t")
	if err == nil {
		t.Errorf("expected an error due to invalid latitude field")
	}

	err = nil
	z, err = ParseTSV("US\t00501\tHoltsville\tNew York\tNY\tSuffolk\t103\t\t\t40.9223\tBadLongitude\t")
	if err == nil {
		t.Errorf("expected an error due to invalid longitude field")
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

	zips, err = LoadTSVFile("BogusFileNameThatHopefullyDoesNotExistOnThisTestSystem")
	if err == nil {
		t.Error(`should have failed to load this file`)
	}
}

func TestDistance(t *testing.T) {
	zips, err := LoadTSVFile("US.txt")
	if err != nil {
		t.Error(err)
		return
	}

	z1 := Find("28115", zips)
	if z1 == nil {
		t.Errorf(`Find("28115", zips) = nil want a valid *Zip`)
		return
	}

	z2 := Find("24450", zips)
	if z2 == nil {
		t.Errorf(`Find("24450", zips) = nil want a valid *Zip`)
		return
	}

	delta := 0.01
	d := Distance(z1, z1)
	if !floatEqual(d, 0.0, delta) {
		t.Errorf("Distance(zip, zip) = %f want 0 +/- %f", d, delta)
	}

	d = Distance(z1, z2)
	want := 170.44
	if !floatEqual(d, want, delta) {
		t.Errorf("Distance(z1, z2) = %f want %f +/- %f", d, want, delta)
	}
}

func TestFind(t *testing.T) {
	zips, err := LoadTSVFile("US.txt")
	if err != nil {
		t.Error(err)
		return
	}

	zip := Find("28115", zips)
	if zip == nil {
		t.Errorf(`Find("28115", zips) = nil want a valid *Zip`)
		return
	}

	if zip.Code != "28115" {
		t.Errorf(`Find("28115", zips) found the wrong zip! zip.Code = %d want "28115"`, zip.Code)
	}

	zip = Find("BogusZip", zips)
	if zip != nil {
		t.Error(`Find("BogusZip", zips) wasn't supposed to find a zip`)
	}
}

func TestFindInRadius(t *testing.T) {
	zips, err := LoadTSVFile("US.txt")
	if err != nil {
		t.Error(err)
		return
	}

	found := FindInRadius("28115", 10, zips)
	want := 13
	if len(found) != want {
		t.Errorf(`FindInRadius("28115", 10, zips) found %d zip codes.  Want %d`, len(found), want)
	}

	zips, err = LoadTSVFile("US.txt")
	if err != nil {
		t.Error(err)
		return
	}

	found = FindInRadius("28115", 10, zips)
	want = 13
	if len(found) != want {
		t.Errorf(`FindInRadius("28115", 10, zips) found %d zip codes.  Want %d`, len(found), want)
	}

	found = FindInRadius("28115", 0, zips)
	want = 1
	if len(found) != want {
		t.Errorf(`FindInRadius("28115", 0, zips) found %d zip codes.  Want %d`, len(found), want)
	}

	found = FindInRadius("BogusZip", 10, zips)
	want = 0
	if len(found) != want {
		t.Errorf(`FindInRadius("BogusZip", 0, zips) found %d zip codes.  Want %d`, len(found), want)
	}
}

func BenchmarkFindInRadius0(b *testing.B)   { benchmarkFindInRadius(0, b) }
func BenchmarkFindInRadius5(b *testing.B)   { benchmarkFindInRadius(5, b) }
func BenchmarkFindInRadius10(b *testing.B)  { benchmarkFindInRadius(10, b) }
func BenchmarkFindInRadius20(b *testing.B)  { benchmarkFindInRadius(20, b) }
func BenchmarkFindInRadius50(b *testing.B)  { benchmarkFindInRadius(50, b) }
func BenchmarkFindInRadius100(b *testing.B) { benchmarkFindInRadius(100, b) }
func BenchmarkFindInRadius500(b *testing.B) { benchmarkFindInRadius(500, b) }

var result []*Zip

func benchmarkFindInRadius(radius float64, b *testing.B) {
	zips, err := LoadTSVFile("US.txt")
	if err != nil {
		b.Error(err)
		return
	}

	var r []*Zip

	for n := 0; n < b.N; n++ {
		r = FindInRadius("28115", radius, zips)
	}

	result = r
}

func floatEqual(a, b, delta float64) bool {
	if (b-delta) <= a && a <= (b+delta) {
		return true
	}
	return false
}
