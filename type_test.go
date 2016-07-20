package eureka_test

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/go-eureka"
)

var _ = Describe("Instance", func() {
	var (
		instanceXml []byte
		instance    = eureka.Instance{
			XMLName:        xml.Name{Local: "instance"},
			Id:             "id",
			HostName:       "host",
			AppName:        "myapp",
			IpAddr:         "1.2.3.4",
			VipAddr:        "vip.address",
			SecureVipAddr:  "secure.vip.address",
			Status:         eureka.StatusUp,
			StatusOverride: eureka.StatusUnknown,
			Port:           80,
			SecurePort:     443,
			HomePageUrl:    "home.page.url",
			StatusPageUrl:  "status.page.url",
			HealthCheckUrl: "health.check.url",
			LeaseInfo: eureka.Lease{
				RenewalInterval:  eureka.Duration(30 * time.Second),
				Duration:         eureka.Duration(90 * time.Second),
				RegistrationTime: eureka.Time(time.Unix(0, 1468519783576*int64(time.Millisecond))),
				LastRenewalTime:  eureka.Time(time.Unix(0, 1468519783577*int64(time.Millisecond))),
				EvictionTime:     eureka.Time(time.Unix(0, 1468519783578*int64(time.Millisecond))),
				ServiceUpTime:    eureka.Time(time.Unix(0, 1468519783579*int64(time.Millisecond))),
			},
			DataCenterInfo: eureka.DataCenter{
				Type: eureka.DataCenterTypePrivate,
				Metadata: eureka.AmazonMetadata{
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
				"b": "two",
				"a": "one",
			},
		}
	)

	BeforeEach(func() {
		var err error
		instanceXml, err = ioutil.ReadFile(filepath.Join("fixtures", "instance.xml"))
		Expect(err).ToNot(HaveOccurred())
		instanceXml = removeIdendation(instanceXml)
	})

	It("can be marshaled to an XML string", func() {
		data, err := xml.Marshal(instance)
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(Equal(instanceXml))
	})

	It("can be unmarshaled from an XML string", func() {
		var actual eureka.Instance
		err := xml.Unmarshal(instanceXml, &actual)
		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(instance))
	})
})