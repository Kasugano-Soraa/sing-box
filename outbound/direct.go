package outbound

import (
	"context"
	"net"
	"net/netip"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/dialer"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"

	"github.com/pires/go-proxyproto"
)

var _ adapter.Outbound = (*Direct)(nil)

type Direct struct {
	myOutboundAdapter
	dialer              N.Dialer
	overrideOption      int
	overrideDestination M.Socksaddr
	proxyProto          uint8
}

func NewDirect(router adapter.Router, logger log.ContextLogger, tag string, options option.DirectOutboundOptions) (*Direct, error) {
	outbound := &Direct{
		myOutboundAdapter: myOutboundAdapter{
			protocol: C.TypeDirect,
			network:  []string{N.NetworkTCP, N.NetworkUDP},
			router:   router,
			logger:   logger,
			tag:      tag,
		},
		dialer:     dialer.NewOutbound(router, options.OutboundDialerOptions),
		proxyProto: options.ProxyProtocol,
	}
	if options.ProxyProtocol > 2 {
		return nil, E.New("invalid proxy protocol option: ", options.ProxyProtocol)
	}
	if options.OverrideAddress != "" && options.OverridePort != 0 {
		outbound.overrideOption = 1
		outbound.overrideDestination = M.ParseSocksaddrHostPort(options.OverrideAddress, options.OverridePort)
	} else if options.OverrideAddress != "" {
		outbound.overrideOption = 2
		outbound.overrideDestination = M.ParseSocksaddrHostPort(options.OverrideAddress, options.OverridePort)
	} else if options.OverridePort != 0 {
		outbound.overrideOption = 3
		outbound.overrideDestination = M.Socksaddr{Port: options.OverridePort}
	}
	return outbound, nil
}

func (h *Direct) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	ctx, metadata := adapter.AppendContext(ctx)
	originDestination := metadata.Destination
	metadata.Outbound = h.tag
	metadata.Destination = destination
	switch h.overrideOption {
	case 1:
		destination = h.overrideDestination
	case 2:
		newDestination := h.overrideDestination
		newDestination.Port = destination.Port
		destination = newDestination
	case 3:
		destination.Port = h.overrideDestination.Port
	}
	network = N.NetworkName(network)
	switch network {
	case N.NetworkTCP:
		h.logger.InfoContext(ctx, "outbound connection to ", destination)
	case N.NetworkUDP:
		h.logger.InfoContext(ctx, "outbound packet connection to ", destination)
	}
	conn, err := h.dialer.DialContext(ctx, network, destination)
	if err != nil {
		return nil, err
	}
	if h.proxyProto > 0 {
		source := metadata.Source
		if !source.IsValid() {
			source = M.SocksaddrFromNet(conn.LocalAddr())
		}
		if originDestination.Addr.Is6() {
			source = M.SocksaddrFrom(netip.AddrFrom16(source.Addr.As16()), source.Port)
		}
		header := proxyproto.HeaderProxyFromAddrs(h.proxyProto, source.TCPAddr(), originDestination.TCPAddr())
		_, err = header.WriteTo(conn)
		if err != nil {
			conn.Close()
			return nil, E.Cause(err, "write proxy protocol header")
		}
	}
	return conn, nil
}

func (h *Direct) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	ctx, metadata := adapter.AppendContext(ctx)
	metadata.Outbound = h.tag
	metadata.Destination = destination
	h.logger.InfoContext(ctx, "outbound packet connection")
	return h.dialer.ListenPacket(ctx, destination)
}

func (h *Direct) NewConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	return NewConnection(ctx, h, conn, metadata)
}

func (h *Direct) NewPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext) error {
	return NewPacketConnection(ctx, h, conn, metadata)
}
