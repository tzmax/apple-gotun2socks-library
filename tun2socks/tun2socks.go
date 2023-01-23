/*
 * @Author: tzmax
 * @Date: 2023-01-22
 * @FilePath: /apple-gotun2socks-library/tun2socks/tun2socks.go
 */

package tun2socks

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

const (
	MTU = 1500
)

type Tun2socksCtl struct {
	TunName string
}

type TunReadWriter interface {
	io.ReadWriteCloser
}

func init() {
	// Apple VPN extensions have a memory limit of 15MB. Conserve memory by increasing garbage
	// collection frequency and returning memory to the OS every minute.
	debug.SetGCPercent(10)
	// TODO: Check if this is still needed in go 1.13, which returns memory to the OS
	// automatically.
	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for range ticker.C {
			debug.FreeOSMemory()
		}
	}()
}

// Tun2socksConnect reads packets from a TUN device and routes it to a socks5 server.
// Returns an Tun2socksCtl instance.
//
// `tunAddr` TUN address.
// `tunWg` TUN Gateway address.
// `tunMask` TUN Masking.
// `tunDns` TUN DNS address.
// `socks5Proxy` socks5 proxy link.
// `isUDPEnabled` indicates whether the tunnel and/or network enable UDP proxying.
//
// Sets an error if the tunnel fails to connect.

func CreateTunConnect(tunAddr, tunWg, tunMask, tunDns, socks5Proxy string, isUDPEnabled bool) (*Tun2socksCtl, error) {
	// Open the tun device.
	if tunDns == "" {
		tunDns = "8.8.8.8,8.8.4.4,1.1.1.1"
	}

	dnsServers := strings.Split(tunDns, ",")
	utunName, tunDev, err := openTunDevice("utun0", tunAddr, tunWg, tunMask, dnsServers, false)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open tun device: %v", err))
	}

	ctl, err := Connect(tunDev, socks5Proxy, isUDPEnabled)
	if err != nil {
		return nil, err
	}

	ctl.TunName = utunName

	return ctl, nil
}

// Tun2socksConnect reads packets from a TUN device and routes it to a socks5 server.
// Returns an Tun2socksCtl instance.
//
// `tunReadWriter` TUN ReadWriter.
// `socks5Proxy` socks5 proxy link.
// `isUDPEnabled` indicates whether the tunnel and/or network enable UDP proxying.
//
// Sets an error if the tunnel fails to connect.

func Connect(tunReadWriter TunReadWriter, socks5Proxy string, isUDPEnabled bool) (*Tun2socksCtl, error) {

	// Setup TCP/IP stack.
	lwipWriter := core.NewLWIPStack().(io.Writer)

	// Register TCP and UDP handlers to handle accepted connections.
	if !strings.Contains(socks5Proxy, "://") {
		socks5Proxy = fmt.Sprintf("socks5://%s", socks5Proxy)
	}
	socksURL, err := url.Parse(socks5Proxy)
	if err != nil {
		return nil, err
	}
	address := socksURL.Host
	if address == "" {
		// Socks5 over UDS
		address = socksURL.Path
	}

	proxyAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid proxy server address: %v", err))
	}
	proxyHost := proxyAddr.IP.String()
	proxyPort := uint16(proxyAddr.Port)

	core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, proxyPort))
	if isUDPEnabled {
		core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, proxyPort, (30 * time.Second)))
	}

	// Register an output callback to write packets output from lwip stack to tun
	// device, output function should be set before input any packets.
	core.RegisterOutputFn(func(data []byte) (int, error) {
		return tunReadWriter.Write(data)
	})

	// Copy packets from tun device to lwip stack, it's the main loop.
	go func() {
		_, err := io.CopyBuffer(lwipWriter, tunReadWriter, make([]byte, MTU))
		if err != nil {
			log.Fatalf("copying data failed: %v", err)
		}
	}()

	return &Tun2socksCtl{}, nil
}
