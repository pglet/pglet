package server

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	trustedCIDRs    []*net.IPNet
	remoteIPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
)

// ClientIP implements a best effort algorithm to return the real client IP.
// It called c.RemoteIP() under the hood, to check if the remote IP is a trusted proxy or not.
// If it's it will then try to parse the headers defined in Engine.RemoteIPHeaders (defaulting to [X-Forwarded-For, X-Real-Ip]).
// If the headers are nots syntactically valid OR the remote IP does not correspong to a trusted proxy,
// the remote IP (coming form Request.RemoteAddr) is returned.
func ClientIP(c *gin.Context) string {
	remoteIP, trusted := RemoteIP(c)
	if remoteIP == nil {
		return ""
	}

	if trusted && remoteIPHeaders != nil {
		for _, headerName := range remoteIPHeaders {
			ip, valid := validateHeader(c.Request.Header.Get(headerName))
			if valid {
				return ip
			}
		}
	}
	return remoteIP.String()
}

// RemoteIP parses the IP from Request.RemoteAddr, normalizes and returns the IP (without the port).
// It also checks if the remoteIP is a trusted proxy or not.
// In order to perform this validation, it will see if the IP is contained within at least one of the CIDR blocks
// defined in Engine.TrustedProxies
func RemoteIP(c *gin.Context) (net.IP, bool) {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if err != nil {
		return nil, false
	}
	remoteIP := net.ParseIP(ip)
	if remoteIP == nil {
		return nil, false
	}

	if trustedCIDRs != nil {
		for _, cidr := range trustedCIDRs {
			if cidr.Contains(remoteIP) {
				return remoteIP, true
			}
		}
	}

	return remoteIP, false
}

func validateHeader(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}
	items := strings.Split(header, ",")
	for i, ipStr := range items {
		ipStr = strings.TrimSpace(ipStr)
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return "", false
		}

		// We need to return the first IP in the list, but,
		// we should not early return since we need to validate that
		// the rest of the header is syntactically valid
		if i == 0 {
			clientIP = ipStr
			valid = true
		}
	}
	return
}

func prepareTrustedCIDRs(engine *gin.Engine) {
	if engine.TrustedProxies == nil {
		return
	}

	trustedCIDRs = make([]*net.IPNet, 0, len(engine.TrustedProxies))
	for _, trustedProxy := range engine.TrustedProxies {
		if !strings.Contains(trustedProxy, "/") {
			ip := parseIP(trustedProxy)
			if ip == nil {
				log.Errorf("error parsing IP: %s", trustedProxy)
				return
			}

			switch len(ip) {
			case net.IPv4len:
				trustedProxy += "/32"
			case net.IPv6len:
				trustedProxy += "/128"
			}
		}
		_, cidrNet, err := net.ParseCIDR(trustedProxy)
		if err != nil {
			log.Errorf("error parsing CIDR: %s", err)
			return
		}
		trustedCIDRs = append(trustedCIDRs, cidrNet)
	}
}

// parseIP parse a string representation of an IP and returns a net.IP with the
// minimum byte representation or nil if input is invalid.
func parseIP(ip string) net.IP {
	parsedIP := net.ParseIP(ip)

	if ipv4 := parsedIP.To4(); ipv4 != nil {
		// return ip in a 4-byte representation
		return ipv4
	}

	// return ip in a 16-byte representation or nil
	return parsedIP
}
