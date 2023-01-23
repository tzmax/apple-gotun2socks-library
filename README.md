# apple-gotun2socks-library

Go package for building [go-tun2socks](https://github.com/eycorsican/go-tun2socks) library for macOS and iOS

## Prerequisites

- macOS host (iOS, macOS)
- make
- Go >= 1.18
- A C compiler (e.g.: clang, gcc)

## Apple (iOS and macOS)

### Set up

- Xcode
- [gomobile](https://pkg.go.dev/golang.org/x/mobile/cmd/gobind) (installed as needed by `make`)

### Build
```
make clean && make apple
```
This will create `build/apple/Tun2socks.xcframework`.

## API
```
#import <Tun2socks/Tun2socks.h>

NSString* tunAddr = @"10.0.0.0";
NSString* tunWg = @"10.0.0.1";
NSString* tunMask = @"255.255.255.0";
NSString* tunDns = @"8.8.8.8,8.8.4.4,1.1.1.1";
NSString* proxyServer = @"socks5://127.0.0.1:1080";

// `tunAddr` TUN address.
// `tunWg` TUN Gateway address.
// `tunMask` TUN Masking.
// `tunDns` TUN DNS address.
// `socks5Proxy` socks5 proxy link.
// `isUDPEnabled` indicates whether the tunnel and/or network enable UDP proxying.

NSError* err;
Tun2socksTun2socksCtl* ctl = Tun2socksCreateTunConnect(tunAddr, tunWg, tunMask, tunDns, socks5ProxyLink, true, &err);
if (err != NULL) {
    NSLog(@"Tun2socksConnect error:  %@\n", err);
}

NSLog(@"tun fd is %@ \n", ctl.tunName);
```

## Contribute

Please refer to the [contribution guide](/CONTRIBUTING.md)

## Support

V2RayXS: GUI for xray-core on macOS [tzmax/V2RayXS](https://github.com/tzmax/V2RayXS)

## Thanks

[eycorsican/go-tun2socks](https://github.com/eycorsican/go-tun2socks)

[songgao/water](https://github.com/songgao/water)
