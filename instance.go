package jolt

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Instance struct {
	XMLName        xml.Name   `xml:"instance"`
	HostName       string     `xml:"hostName"`
	App            string     `xml:"app"`
	IpAddr         string     `xml:"ipAddr"`
	VipAddr        string     `xml:"vipAddress"`
	SecureVipAddr  string     `xml:"secureVipAddress"`
	Status         status     `xml:"status"`
	Port           int        `xml:"port"`
	SecurePort     int        `xml:"securePort"`
	HomePageUrl    string     `xml:"homePageUrl"`
	StatusPageUrl  string     `xml:"statusPageUrl"`
	HealthCheckUrl string     `xml:"healthCheckUrl"`
	DatacenterInfo Datacenter `xml:"dataCenterInfo"`
	LeaseInfo      Lease      `xml:"leaseInfo"`
	Metadata       Metadata   `xml:"metadata"`
}

type Lease struct {
	EvictionDurationInSecs int `xml:"evictionDurationInSecs"`
}

type Datacenter struct {
	Name     datacenterName `xml:"name"`
	Metadata amazonMetadata `xml:"metadata"`
}

type amazonMetadata struct {
	Hostname         string `xml:"hostname'`
	PublicHostname   string `xml:"public-hostname"`
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

type datacenterName uint8

const (
	DatacenterPrivate datacenterName = iota
	DatacenterAmazon
)

var datacenterNames = []string{
	"MyOwn",
	"Amazon",
}

func (dn datacenterName) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if int(dn) >= len(datacenterNames) {
		return fmt.Errorf("Unknown datacenter code: %d", dn)
	}
	return e.EncodeElement(datacenterNames[dn], start)
}

func (dn *datacenterName) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}

	for i, n := range datacenterNames {
		if n == str {
			*dn = datacenterName(i)
			return nil
		}
	}

	return fmt.Errorf("Unknown datacenter name: %s", str)
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

	for key, value := range m {
		if err := e.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: key}}); err != nil {
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
