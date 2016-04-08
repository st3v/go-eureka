package jolt_test

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/jolt"
)

var _ = Describe("Instance", func() {
	var (
		instanceXml []byte
		instance    = jolt.Instance{
			XMLName:        xml.Name{Local: "instance"},
			HostName:       "host",
			App:            "myapp",
			IpAddr:         "1.2.3.4",
			VipAddr:        "vip.address",
			SecureVipAddr:  "secure.vip.address",
			Status:         jolt.StatusOutOfService,
			Port:           80,
			SecurePort:     443,
			HomePageUrl:    "home.page.url",
			StatusPageUrl:  "status.page.url",
			HealthCheckUrl: "health.check.url",
			LeaseInfo: jolt.Lease{
				EvictionDurationInSecs: 123,
			},
			DataCenterInfo: jolt.DataCenter{
				Type: jolt.DataCenterTypePrivate,
				Metadata: jolt.AmazonMetadata{
					Hostname:         "dchost",
					PublicHostName:   "dc.public.host",
					LocalHostName:    "dc.local.host",
					PublicIpv4:       "1.2.3.5",
					LocalIpv4:        "1.2.3.6",
					AvailabilityZone: "az",
					InstanceId:       "instance.id",
					InstanceType:     "instance.type",
					AmiId:            "ami.id",
					AmiLaunchIndex:   "ami.launch.index",
					AmiManifestPath:  "ami.manifest.path",
				},
			},
			Metadata: map[string]string{
				"foo": "one",
				"bar": "two",
			},
		}
	)

	BeforeEach(func() {
		var err error
		instanceXml, err = ioutil.ReadFile(filepath.Join("fixtures", "instance.xml"))
		Expect(err).ToNot(HaveOccurred())
	})

	It("can be marshaled to an XML string", func() {
		data, err := xml.MarshalIndent(instance, "", "    ")
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(Equal(instanceXml))
	})

	It("can be unmarshaled from an XML string", func() {
		var actual jolt.Instance
		err := xml.Unmarshal(instanceXml, &actual)
		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(instance))
	})
})
