package eureka

import (
	"encoding/xml"
	"strings"
	"time"
)

type Instance struct {
	XMLName        xml.Name   `xml:"instance"`
	ID             string     `xml:"instanceId"`
	HostName       string     `xml:"hostName"`
	AppName        string     `xml:"app"`
	IPAddr         string     `xml:"ipAddr"`
	VIPAddr        string     `xml:"vipAddress"`
	SecureVIPAddr  string     `xml:"secureVipAddress"`
	Status         Status     `xml:"status"`
	StatusOverride Status     `xml:"overriddenstatus"`
	Port           Port       `xml:"port"`
	SecurePort     Port       `xml:"securePort"`
	HomePageURL    string     `xml:"homePageUrl"`
	StatusPageURL  string     `xml:"statusPageUrl"`
	HealthCheckURL string     `xml:"healthCheckUrl"`
	DataCenterInfo DataCenter `xml:"dataCenterInfo"`
	LeaseInfo      Lease      `xml:"leaseInfo"`
	Metadata       Metadata   `xml:"metadata"`
}

// Equals checks if two instances are the same. Does not compare LeaseInfo.
func (i *Instance) Equals(other *Instance) bool {
	return i.ID == other.ID &&
		i.HostName == other.HostName &&
		strings.ToUpper(i.AppName) == strings.ToUpper(other.AppName) &&
		i.IPAddr == other.IPAddr &&
		i.VIPAddr == other.VIPAddr &&
		i.SecureVIPAddr == other.SecureVIPAddr &&
		i.Status == other.Status &&
		i.StatusOverride == other.StatusOverride &&
		i.Port == other.Port &&
		i.SecurePort == other.SecurePort &&
		i.HomePageURL == other.HomePageURL &&
		i.StatusPageURL == other.StatusPageURL &&
		i.HealthCheckURL == other.HealthCheckURL &&
		i.DataCenterInfo == other.DataCenterInfo &&
		i.Metadata.Equals(other.Metadata)
}

type Port uint16

type Status uint8

const (
	StatusUp Status = iota
	StatusDown
	StatusStarting
	StatusOutOfService
	StatusUnknown
)

type DataCenter struct {
	Type     DataCenterType `xml:"name"`
	Metadata AmazonMetadata `xml:"metadata"`
}

type DataCenterType uint8

const (
	DataCenterTypePrivate DataCenterType = iota
	DataCenterTypeAmazon
)

type AmazonMetadata struct {
	HostName         string `xml:"hostname"`
	PublicHostName   string `xml:"public-hostname"`
	LocalHostName    string `xml:"local-hostname"`
	PublicIPV4       string `xml:"public-ipv4"`
	LocalIPV4        string `xml:"local-ipv4"`
	AvailabilityZone string `xml:"availability-zone"`
	InstanceID       string `xml:"instance-id"`
	InstanceType     string `xml:"instance-type"`
	AmiID            string `xml:"ami-id"`
	AmiLaunchIndex   string `xml:"ami-launch-index"`
	AmiManifestPath  string `xml:"ami-manifest-path"`
}

type Lease struct {
	RenewalInterval  Duration `xml:"renewalIntervalInSecs"`
	Duration         Duration `xml:"durationInSecs"`
	RegistrationTime Time     `xml:"registrationTimestamp"`
	LastRenewalTime  Time     `xml:"lastRenewalTimestamp"`
	EvictionTime     Time     `xml:"evictionTimestamp"`
	ServiceUpTime    Time     `xml:"serviceUpTimestamp"`
}

type Duration time.Duration

type Time time.Time

type Metadata map[string]string

func (m Metadata) Equals(other Metadata) bool {
	if len(m) != len(other) {
		return false
	}

	for k, v := range m {
		if v != other[k] {
			return false
		}
	}

	return true
}

type App struct {
	XMLName   xml.Name    `xml:"application"`
	Name      string      `xml:"name"`
	Instances []*Instance `xml:"instance"`
}

type Registry struct {
	XMLName      xml.Name `xml:"applications"`
	VersionDelta int      `xml:"versions__delta"`
	Hashcode     string   `xml:"apps__hashcode"`
	Apps         []*App   `xml:"application"`
}
