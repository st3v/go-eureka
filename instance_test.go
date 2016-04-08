package jolt_test

import (
	"encoding/xml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/jolt"
)

var _ = Describe("Instance", func() {
	var (
		instanceXml = []byte("<instance><hostName>hostname</hostName><app>app</app><ipAddr>1.2.3.4</ipAddr><vipAddress>hostname.domain</vipAddress><secureVipAddress></secureVipAddress><status>UP</status><port>0</port><securePort>0</securePort><homePageUrl></homePageUrl><statusPageUrl></statusPageUrl><healthCheckUrl></healthCheckUrl><dataCenterInfo><name>MyOwn</name><metadata><Hostname></Hostname><public-hostname></public-hostname><local-hostname></local-hostname><PublicIpv4></PublicIpv4><local-ipv4></local-ipv4><availability-zone></availability-zone><instance-id></instance-id><instance-type></instance-type><ami-id></ami-id><ami-launch-index></ami-launch-index><ami-manifest-path></ami-manifest-path></metadata></dataCenterInfo><leaseInfo><evictionDurationInSecs>0</evictionDurationInSecs></leaseInfo><metadata><foo>one</foo><bar>two</bar></metadata></instance>")
	)

	It("can be marshaled to an XML string", func() {
		instance := jolt.Instance{
			HostName: "hostname",
			App:      "app",
			IpAddr:   "1.2.3.4",
			VipAddr:  "hostname.domain",
			Metadata: map[string]string{"foo": "one", "bar": "two"},
		}

		xmlStr, err := xml.Marshal(instance)

		Expect(err).ToNot(HaveOccurred())
		Expect(xmlStr).To(Equal(instanceXml))

	})

	It("can be unmarshaled from an XML string", func() {
		var instance jolt.Instance
		err := xml.Unmarshal(instanceXml, &instance)

		Expect(err).ToNot(HaveOccurred())
		Expect(instance.HostName).To(Equal("hostname"))
	})
})
