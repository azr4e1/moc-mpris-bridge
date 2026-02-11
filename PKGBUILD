# Maintainer: Your Name <your@email.com>
pkgname=moc-mpris-bridge
pkgver=1.0.0
pkgrel=1
pkgdesc="MPRIS bridge for Music on Console"
arch=('x86_64' 'aarch64')
url="https://github.com/azr4e1/moc-mpris-bridge"
license=('MIT')  # adjust to your license
depends=('moc-pulse')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')  # run updpkgsums after first build

build() {
    cd "$pkgname-$pkgver"
    export CGO_CPPFLAGS="${CPPFLAGS}"
    export CGO_CFLAGS="${CFLAGS}"
    export CGO_CXXFLAGS="${CXXFLAGS}"
    export CGO_LDFLAGS="${LDFLAGS}"
    export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
    go build -o "$pkgname"
}

package() {
    cd "$pkgname-$pkgver"
    install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
    
    # optional: install systemd user service
    install -Dm644 "$pkgname.service" \
        "$pkgdir/usr/lib/systemd/user/$pkgname.service"
}
