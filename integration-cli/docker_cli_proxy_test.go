package main

import (
	"net"
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

func (s *DockerSuite) TestCLIProxyDisableProxyUnixSock(c *check.C) {
	testRequires(c, DaemonIsLinux, SameHostDaemon)

	icmd.RunCmd(icmd.Cmd{
		Command: []string{dockerBinary, "info"},
		Env:     appendBaseEnv(false, "HTTP_PROXY=http://127.0.0.1:9999"),
	}).Assert(c, icmd.Success)
}

// Can't use localhost here since go has a special case to not use proxy if connecting to localhost
// See https://golang.org/pkg/net/http/#ProxyFromEnvironment
func (s *DockerDaemonSuite) TestCLIProxyProxyTCPSock(c *check.C) {

	c.Skip("Pending balenaEngine compatibility investigation")

	testRequires(c, SameHostDaemon)
	// get the IP to use to connect since we can't use localhost
	addrs, err := net.InterfaceAddrs()
	c.Assert(err, checker.IsNil)
	var ip string
	for _, addr := range addrs {
		sAddr := addr.String()
		if !strings.Contains(sAddr, "127.0.0.1") {
			addrArr := strings.Split(sAddr, "/")
			ip = addrArr[0]
			break
		}
	}

	c.Assert(ip, checker.Not(checker.Equals), "")

	s.d.Start(c, "-H", "tcp://"+ip+":2375")

	icmd.RunCmd(icmd.Cmd{
		Command: []string{dockerBinary, "info"},
		Env:     []string{"DOCKER_HOST=tcp://" + ip + ":2375", "HTTP_PROXY=127.0.0.1:9999"},
	}).Assert(c, icmd.Expected{Error: "exit status 1", ExitCode: 1})
	// Test with no_proxy
	icmd.RunCmd(icmd.Cmd{
		Command: []string{dockerBinary, "info"},
		Env:     []string{"DOCKER_HOST=tcp://" + ip + ":2375", "HTTP_PROXY=127.0.0.1:9999", "NO_PROXY=" + ip},
	}).Assert(c, icmd.Success)
}
