package publicip

import (
	"context"
	"net"
)

var openDnsResolver = &net.Resolver{
	// This ensures our Dial function will be used.
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			// I think we want to disable switching to other networks, since this DNS response
			// depends on the network used.
			FallbackDelay: -1,
		}
		// Go DNS resolution doesn't bother to use a particular network, but we need to reapply our
		// constraint.
		suffix := ctx.Value(dialNetworkSuffixVar)
		if suffix != nil {
			network += suffix.(string)
		}
		return d.DialContext(ctx, network, "resolver1.opendns.com:53")
	},
}

var dialNetworkSuffixVar = &struct{}{}

func withDialNetworkSuffix(ctx context.Context, suffix string) context.Context {
	return context.WithValue(ctx, dialNetworkSuffixVar, suffix)
}

func GetAll(ctx context.Context) ([]net.IPAddr, error) {
	// We know this passes "ip" to the LookupIP, and includes Zones in the return.
	return openDnsResolver.LookupIPAddr(ctx, "myip.opendns.com")
}

// Network should be one of "ip", "ip4", or "ip6".
func Get(ctx context.Context, network string) ([]net.IP, error) {
	return openDnsResolver.LookupIP(ctx, network, "myip.opendns.com")
}

func Get4(ctx context.Context) (net.IP, error) {
	all, err := Get(withDialNetworkSuffix(ctx, "4"), "ip4")
	if err != nil {
		return nil, err
	}
	return all[0], nil
}

func Get6(ctx context.Context) (net.IP, error) {
	all, err := Get(withDialNetworkSuffix(ctx, "6"), "ip6")
	if err != nil {
		return nil, err
	}
	return all[0], nil
}
