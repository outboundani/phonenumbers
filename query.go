package phonenumbers

import (
	"strings"
	"unicode"
)

// Record contains all the data we have or have inferred for a phone number by parsing it and deconstructing it
type Record struct {

	// International format is the usual +countrycode<number> format
	International string

	// NType is the type of number, e.g. fixed line, mobile, toll free etc
	NType PhoneNumberType

	// Geocode is usually the town or city where the number is registered, including the Admin1 level such as state or province
	Geocode string

	// Timezone is the timezone of the number
	Timezone string

	// Region is usually the country, though some countries have more than one region
	Region string

	// The language used to return the geocode and region etc
	Language string

	// Country code is the + prefix for a number
	CountryCode string

	// NDC is the National Destination Code, which is the area code for the number in the US for instance but
	// varies a little nation to nation
	NDC string

	// LDC is the Local Destination Code, which is the local exchange code for the number in the US for instance but
	// is not always present, for instance UK numbers only have teh Local field
	LDC string

	// Local is the local part of the number, which is the subscriber number (last three digits) in the US for instance
	// but is the whole of hte non NDC part of the number for many countries such as the UK
	Local string

	// Carrier is the name of the carrier for the number if we can work this out. This is not always possible
	Carrier string

	// E164 is the number in E164 format, which is the international format without the + prefix and is often used for
	// storage keys
	E164 string

	// Valid indicates whether the phone parser thinks that this is a valid number for the give region or not
	Valid bool
}

// ParsePhone takes a phone number in national or international format and returns a Record containing all the data we have
// about this phone number. This includes the international format, the type of number, the geocoded location, the timezone
// etc,
//
// Note that the region parameter is usually the country code in 2 character format, but some countries have more than one
// region; it defaults to "US".
//
// The lang parameter is the language to use for the geocoding and carrier lookup. If this is not specified, then it defaults to "en"
//
//goland:noinspection GoUnusedExportedFunction
func ParsePhone(number string, region string, lang string) (Record, error) {

	var rec Record
	if region == "" {
		region = "US"
	}
	if lang == "" {
		lang = "en"
	}

	n, err := Parse(number, region)
	if err != nil {
		return rec, err
	}
	rec.Valid = IsValidNumber(n)
	rec.International = Format(n, INTERNATIONAL)
	rec.NType = GetNumberType(n)

	gc, err := GetGeocodingForNumber(n, lang)
	if err == nil {
		rec.Geocode = gc
	}

	tz, err := GetTimezonesForNumber(n)
	if err == nil {
		rec.Timezone = tz[0]
	}

	rec.Region = GetRegionCodeForNumber(n)
	rec.E164 = Format(n, E164)

	parts := SplitPhone(rec.International)
	rec.CountryCode = parts[0]
	rec.NDC = parts[1]
	switch len(parts) {
	case 3: // UK etc
		rec.Local = parts[2]
	case 4: // USA, Canada, Mexico et al
		rec.LDC = parts[2]
		rec.Local = parts[3]
	}
	cr, err := GetCarrierForNumber(n, lang)
	if err == nil {
		rec.Carrier = cr
	}
	return rec, nil
}

// SplitPhone takes a string in international format and returns the separate parts of the
// number as a slice of strings. Because each country's international formatting and number of
// constituent parts is different, we have this customer function as we need to know each of the
// parts for classification fields.
//
// For example, take the US number:
//
//	+1 212-555-1212.
//
// This would be split into 4 parts. The first part is the country code, the second part is the geographical "area code",
// more accurately referred to as the National Destination Code of NDC, the third part is the local code and the fourth
// and last part is the subscriber number. The parts 2, 3, 4 together make up the national (significant) number.
//
// However, Taiwan's international format uses spaces rather than hyphens to separate the parts, so the number formatted
// internationally would look like:
//
//	+886 255 555 1212
//
// However, it still has four parts with the same meaning.
//
// UK numbers, despite the UK inventing everything, have been fucked up by a quango rather than sound principles, so in
// international format, a number may look like:
//
//	+44 7911 123456
//
// Which means that there is an NDC of 7911 but no local exchange code because the NDC already narrows it down by
// having 4 digits. They just had to be different.
//
// Hence, this function splits on anything that is not numeric (+ is included as numeric).
func SplitPhone(phone string) []string {

	// We are just wrapping the built in Go functionality for now, but doing this in one place
	// means that if we need something more complicated later, then it is easy to replace.
	//
	parts := strings.FieldsFunc(phone, func(r rune) bool {
		switch {
		case r == '+':
			return false
		case unicode.IsNumber(r):
			return false
		default:
			return true
		}
	})
	return parts
}
