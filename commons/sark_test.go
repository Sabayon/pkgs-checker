/*

Copyright (C) 2017-2019  Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package commons_test

import (
	"fmt"

	. "github.com/Sabayon/pkgs-checker/commons"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SarkConfig", func() {

	Describe("NewSarkConfigFromString", func() {

		conf := `---
injector:
  filter:
    type: "whitelist"
    rules:
      - categories:
        - "dev-libs"
        - "www-servers"
`
		sark, err := NewSarkConfigFromString(nil, conf)

		out, _ := sark.ToString()
		fmt.Println(fmt.Sprintf("SARK %s", out))

		Context("Check processing phase", func() {
			It("Check attributes", func() {
				Expect(err).Should(BeNil())
			})
		})

		Context("Check filter", func() {
			It("Check attributes", func() {
				Expect(sark.Injector.Filter.FilterType).To(Equal("whitelist"))
			})
		})

	})

	Describe("NewSarkConfigFromString2", func() {

		conf := `---
build:
  equo:
    dependency_install:
      enable: 1
      install_atoms: 1
      # It seems that sometime dependencies calculation
      # doesn't catch all packages. I will investigate
      # split_install: 1
    package:

      unmask:
        - dev-util/xdelta
        - dev-util/ostree
        - net-misc/cni-plugins
        - net-firewall/conntrack-tools
        - app-emulation/lxd
        - app-emulation/cri-o
        - app-emulation/docker-swarm
        - app-emulation/docker-registry
        - app-emulation/docker-compose
        - dev-python/dockerpty
        - dev-python/docker-pycreds
        - dev-python/docker-py
        - dev-python/websocket-client
        - net-libs/libnetfilter_log

      install:
        - dev-lang/swig

  qa_checks: 0
  emerge:
    # Set sabayon-arm profile
    profile: "sabayon-distro:default/linux/arm/17.0/armv7a/sabayon"

    # Install each package separately
    split_install: 1
    preserved_rebuild: 1
    jobs: 3

    default_args: --accept-properties=-interactive -t --complete-graph --newuse --deep --with-bdeps=y --buildpkg --update

    features: assume-digests binpkg-logs -userpriv config-protect-if-modified distlocks ebuild-locks fixlafiles merge-sync parallel-fetch preserve-libs protect-owned sandbox sfperms splitdebug strict

  target:
    - net-dns/dnsmasq
    - net-dns/bind-tools
    - net-dns/bind
    - net-nds/rpcbind
    - dev-util/ostree
    - app-emulation/containerd
    - dev-libs/libnl
    - net-libs/libnetfilter_cthelper
    - net-libs/libnetfilter_cttimeout
    - net-libs/libnetfilter_queue
    - net-libs/libnetfilter_conntrack
    - net-libs/libnetfilter_acct
    - net-libs/libnetfilter_log
    - net-firewall/conntrack-tools
    - net-misc/socat
    - net-misc/cni-plugins
    - app-emulation/cri-o
`
		sark, err := NewSarkConfigFromString(nil, conf)

		out, _ := sark.ToString()
		fmt.Println(fmt.Sprintf("SARK %s", out))

		Context("Check processing phase", func() {
			It("Check attributes", func() {
				Expect(err).Should(BeNil())
			})
		})

		Context("Check filter", func() {
			It("Check attributes", func() {
				Expect(sark.Build.Equo.Packages.Install).To(Equal([]string{"dev-lang/swig"}))
			})
		})

	})

	Describe("NewSarkConfigFromFile", func() {

		sark, err := NewSarkConfigFromFile(nil, "../tests/sark/inject_example1.yaml")

		Context("Check processing phase", func() {
			It("Check attributes", func() {
				Expect(err).Should(BeNil())
			})
		})

		Context("Check filter", func() {
			It("Check attributes", func() {
				Expect(sark.Injector.Filter.FilterType).To(Equal("blacklist"))
			})
		})

	})

})
