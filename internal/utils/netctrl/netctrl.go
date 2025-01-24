package netctrl

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const HeaderIPKey = "X-Real-IP"

type IPControl struct {
	log      *zap.Logger
	whiteNet net.IPNet
}

func NewIPControl(ipnet string, log *zap.Logger) (*IPControl, error) {
	_, net, err := net.ParseCIDR(ipnet)
	if err != nil {
		return nil, fmt.Errorf("error parsing CIDR: %w", err)
	}
	return &IPControl{
		whiteNet: *net,
		log:      log,
	}, nil
}

func (i *IPControl) IsIPAllowed(ip net.IP) bool {
	return i.whiteNet.Contains(ip)
}

func (i *IPControl) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if is := ctx.Request.Header.Get(HeaderIPKey); is != "" {
			if ip := net.ParseIP(is); ip != nil {
				if i.whiteNet.Contains(ip) {
					ctx.Next()
				} else {
					_ = ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("ip %s is not in trusted net", is))
				}
			} else {
				_ = ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("bad ip value %s in header %s", is, HeaderIPKey))
			}
		} else {
			_ = ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("ip expected in header %s", HeaderIPKey))
		}
	}
}

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, fmt.Errorf("error on dial: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if localAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return localAddr.IP, nil
	}

	return nil, errors.New("failed reading local ip")
}
