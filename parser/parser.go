package parser

import (
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-marc/fields"
	"strconv"
	"strings"
)

// please for to be replacing me with interface{} thingies
// defined in bbox.go

type Point struct {
	Latitude  float64
	Longitude float64
}

func (p *Point) String() string {
	return fmt.Sprintf("%0.6f %0.6f", p.Latitude, p.Longitude)
}

type BoundingBox struct {
	SW Point
	NE Point
}

func (b *BoundingBox) String() string {
	return fmt.Sprintf("%s %s", b.SW.String(), b.NE.String())
}

type Parser struct {
	Scheme    string
	Order     string
	Separator string
}

func NewParser() (*Parser, error) {

	p := Parser{
		Scheme:    "nsew",
		Order:     "latlon",
		Separator: ",",
	}

	return &p, nil
}

func (p *Parser) Parse(bbox string) (*BoundingBox, error) {

	switch p.Scheme {

	case "swne":
		return p.ParseCorners(bbox)
	case "nsew":
		return p.ParseSides(bbox)
	case "marc":
		return p.ParseMARC(bbox)
	default:
		return nil, errors.New("Invalid or unsupported bounding box scheme")
	}
}

func (p *Parser) ParseSides(bbox string) (*BoundingBox, error) {

	var str_swlat string
	var str_swlon string
	var str_nelat string
	var str_nelon string

	var err error

	parts := strings.Split(bbox, p.Separator)

	if len(parts) != 4 {
		return nil, errors.New("Invalid bounding box")
	}

	switch p.Scheme {
	case "nsew":
		str_swlat = parts[1]
		str_swlon = parts[3]
		str_nelat = parts[0]
		str_nelon = parts[2]
	default:
		return nil, errors.New("Unsupported or invalid scheme")
	}

	bb, err := p.ParseBbox(str_swlat, str_swlon, str_nelat, str_nelon)

	if err != nil {
		return nil, err
	}

	_, err = p.ValidateBbox(bb)

	if err != nil {
		return nil, err
	}

	return bb, nil
}

func (p *Parser) ParseCorners(bbox string) (*BoundingBox, error) {

	var str_swlat string
	var str_swlon string
	var str_nelat string
	var str_nelon string
	var err error

	parts := strings.Split(bbox, p.Separator)

	if len(parts) != 4 {
		return nil, errors.New("Invalid bounding box")
	}

	switch p.Scheme {

	case "swne":

		if p.Order == "latlon" {

			str_swlat = parts[0]
			str_swlon = parts[1]
			str_nelat = parts[2]
			str_nelon = parts[3]

		} else if p.Order == "lonlat" {

			str_swlat = parts[0]
			str_swlon = parts[1]
			str_nelat = parts[2]
			str_nelon = parts[3]

		} else {
			return nil, errors.New("Invalid ordering")
		}

	default:
		return nil, errors.New("Invalid or unsupported scheme")
	}

	bb, err := p.ParseBbox(str_swlat, str_swlon, str_nelat, str_nelon)

	if err != nil {
		return nil, err
	}

	_, err = p.ValidateBbox(bb)

	if err != nil {
		return nil, err
	}

	return bb, nil
}

func (p *Parser) ParseMARC(bbox string) (*BoundingBox, error) {

	parsed, err := fields.Parse034(bbox)

	if err != nil {
		return nil, errors.New("Invalid 034 MARC string")
	}

	_bb, err := parsed.BoundingBox()

	if err != nil {
		return nil, errors.New("Failed to derive bounding box from 034 MARC string")
	}

	// this is to account for the fact that we don't have an interface{} thingy
	// to share across packages yet... (20170220/thisisaaronland)

	sw := Point{
		Latitude:  _bb.SW.Latitude,
		Longitude: _bb.SW.Longitude,
	}

	ne := Point{
		Latitude:  _bb.NE.Latitude,
		Longitude: _bb.NE.Longitude,
	}

	bb := &BoundingBox{
		SW: sw,
		NE: ne,
	}

	_, err = p.ValidateBbox(bb)

	if err != nil {
		return nil, err
	}

	return bb, nil
}

func (p *Parser) ParseBbox(str_swlat string, str_swlon string, str_nelat string, str_nelon string) (*BoundingBox, error) {

	var swlat float64
	var swlon float64
	var nelat float64
	var nelon float64
	var err error

	swlat, err = strconv.ParseFloat(str_swlat, 64)

	if err != nil {
		return nil, errors.New("Invalid SW latitude parameter")
	}

	swlon, err = strconv.ParseFloat(str_swlon, 64)

	if err != nil {
		return nil, errors.New("Invalid SW longitude parameter")
	}

	nelat, err = strconv.ParseFloat(str_nelat, 64)

	if err != nil {
		return nil, errors.New("Invalid NE latitude parameter")
	}

	nelon, err = strconv.ParseFloat(str_nelon, 64)

	if err != nil {
		return nil, errors.New("Invalid NE longitude parameter")
	}

	sw := Point{
		Latitude:  swlat,
		Longitude: swlon,
	}

	ne := Point{
		Latitude:  nelat,
		Longitude: nelon,
	}

	bbox := BoundingBox{
		SW: sw,
		NE: ne,
	}

	return &bbox, nil
}

func (p *Parser) ValidateBbox(bbox *BoundingBox) (bool, error) {

	swlat := bbox.SW.Latitude
	swlon := bbox.SW.Longitude
	nelat := bbox.NE.Latitude
	nelon := bbox.NE.Longitude

	if swlat > 90.0 || swlat < -90.0 {
		return false, errors.New("E_IMPOSSIBLE_LATITUDE (SW)")
	}

	if nelat > 90.0 || nelat < -90.0 {
		return false, errors.New("E_IMPOSSIBLE_LATITUDE (NE)")
	}

	if swlon > 180.0 || swlon < -180.0 {
		return false, errors.New("E_IMPOSSIBLE_LONGITUDE (SW)")
	}

	if nelon > 180.0 || nelon < -180.0 {
		return false, errors.New("E_IMPOSSIBLE_LONGITUDE (ne)")
	}

	if swlat > nelat {
		return false, errors.New("E_INVALID_LATITUDE (SW > NE)")
	}

	if swlon > nelon {
		return false, errors.New("E_INVALID_LATITUDE (SW > NE)")
	}

	return true, nil
}
