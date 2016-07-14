package eureka

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"time"
)

var dataCenterTypes = []string{
	"MyOwn",
	"Amazon",
}

func (dct DataCenterType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if int(dct) >= len(dataCenterTypes) {
		return fmt.Errorf("Unknown datacenter type code: %d", dct)
	}
	return e.EncodeElement(dataCenterTypes[dct], start)
}

func (dct *DataCenterType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}

	for i, n := range dataCenterTypes {
		if n == str {
			*dct = DataCenterType(i)
			return nil
		}
	}

	return fmt.Errorf("Unknown datacenter type: %s", str)
}

func (m Metadata) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for _, key := range keys {
		if err := e.EncodeElement(m[key], xml.StartElement{Name: xml.Name{Local: key}}); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func (m *Metadata) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	aux := make(map[string]string)

	for {
		token, err := d.Token()
		if err == io.EOF {
			break
		}

		if start, ok := token.(xml.StartElement); ok {
			var str string
			if err := d.DecodeElement(&str, &start); err != nil {
				return err
			}
			aux[start.Name.Local] = str
		}
	}

	*m = aux

	return nil
}

type portXML struct {
	Value   Port `xml:",chardata"`
	Enabled bool `xml:"enabled,attr"`
}

func (p Port) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(portXML{p, p != 0}, start)
}

func (p *Port) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var aux portXML
	if err := d.DecodeElement(&aux, &start); err != nil {
		return err
	}

	*p = aux.Value

	return nil
}

func (d Duration) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(int64(time.Duration(d).Seconds()), start)
}

func (d *Duration) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var seconds int64
	if err := decoder.DecodeElement(&seconds, &start); err != nil {
		return err
	}

	*d = Duration(time.Duration(seconds) * time.Second)

	return nil
}

func (t Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	epoch := int64(time.Time(t).UnixNano() / int64(time.Millisecond))
	return e.EncodeElement(epoch, start)
}

func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var epoch int64
	if err := d.DecodeElement(&epoch, &start); err != nil {
		return err
	}

	*t = Time(time.Unix(0, epoch*int64(time.Millisecond)))

	return nil
}

var statusNames = []string{
	"UP",
	"DOWN",
	"STARTING",
	"OUT_OF_SERVICE",
	"UNKNOWN",
}

func (s Status) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if int(s) >= len(statusNames) {
		return fmt.Errorf("Unknown status code: %d", s)
	}
	return e.EncodeElement(statusNames[s], start)
}

func (s *Status) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}

	for i, n := range statusNames {
		if n == str {
			*s = Status(i)
			return nil
		}
	}

	return fmt.Errorf("Unknown status: %s", str)
}
