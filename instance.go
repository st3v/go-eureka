package eureka

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
)

type App struct {
	XMLName   xml.Name   `xml:"application"`
	Name      string     `xml:"name"`
	Instances []Instance `xml:"instance"`
}

type Instance struct {
	XMLName        xml.Name   `xml:"instance"`
	Id             string     `xml:"instanceId"`
	HostName       string     `xml:"hostName"`
	AppName        string     `xml:"app"`
	IpAddr         string     `xml:"ipAddr"`
	VipAddr        string     `xml:"vipAddress"`
	SecureVipAddr  string     `xml:"secureVipAddress"`
	Status         status     `xml:"status"`
	Port           int        `xml:"port"`
	SecurePort     int        `xml:"securePort"`
	HomePageUrl    string     `xml:"homePageUrl"`
	StatusPageUrl  string     `xml:"statusPageUrl"`
	HealthCheckUrl string     `xml:"healthCheckUrl"`
	DataCenterInfo DataCenter `xml:"dataCenterInfo"`
	LeaseInfo      Lease      `xml:"leaseInfo"`
	Metadata       Metadata   `xml:"metadata"`
}

type Lease struct {
	EvictionDurationInSecs int `xml:"evictionDurationInSecs"`
}

type DataCenter struct {
	Type     dataCenterType `xml:"name"`
	Metadata AmazonMetadata `xml:"metadata"`
}

type AmazonMetadata struct {
	Hostname         string `xml:"hostname'`
	PublicHostName   string `xml:"public-hostname"`
	LocalHostName    string `xml:"local-hostname"`
	PublicIpv4       string `xml:"public-ipv4'`
	LocalIpv4        string `xml:"local-ipv4"`
	AvailabilityZone string `xml:"availability-zone"`
	InstanceId       string `xml:"instance-id"`
	InstanceType     string `xml:"instance-type"`
	AmiId            string `xml:"ami-id"`
	AmiLaunchIndex   string `xml:"ami-launch-index"`
	AmiManifestPath  string `xml:"ami-manifest-path"`
}

type dataCenterType uint8

const (
	DataCenterTypePrivate dataCenterType = iota
	DataCenterTypeAmazon
)

var dataCenterTypes = []string{
	"MyOwn",
	"Amazon",
}

func (dn dataCenterType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if int(dn) >= len(dataCenterTypes) {
		return fmt.Errorf("Unknown datacenter type code: %d", dn)
	}
	return e.EncodeElement(dataCenterTypes[dn], start)
}

func (dn *dataCenterType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}

	for i, n := range dataCenterTypes {
		if n == str {
			*dn = dataCenterType(i)
			return nil
		}
	}

	return fmt.Errorf("Unknown datacenter type: %s", str)
}

type status uint8

const (
	StatusUp status = iota
	StatusDown
	StatusStarting
	StatusOutOfService
	StatusUnknown
)

var statusNames = []string{
	"UP",
	"DOWN",
	"STARTING",
	"OUT_OF_SERVICE",
	"UNKNOWN",
}

func (s status) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if int(s) >= len(statusNames) {
		return fmt.Errorf("Unknown status code: %d", s)
	}
	return e.EncodeElement(statusNames[s], start)
}

func (s *status) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}

	for i, n := range statusNames {
		if n == str {
			*s = status(i)
			return nil
		}
	}

	return fmt.Errorf("Unknown status: %s", str)
}

type Metadata map[string]string

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
