# Maintainer: Jay C <cuevas0212@gmail.com>
pkgname="gogurt"
pkgver=r54.03c1e7e
pkgrel=1
pkgdesc="A cli app for common things I do"
arch=('i686' 'pentium4' 'x86_64' 'arm' 'armv6h' 'armv7h' 'aarch64')
url="https://github.com/CavemanJay/gogurt"
license=('GPL')
depends=()
makedepends=('go' 'git')
optdepends=('ffmpeg: For converting yt videos to audio'
            'p7zip: For archiving purposes'
            'youtube-dl: For downloading yt videos')
conflicts=('gogurt')
provides=('gogurt')
source=("gogurt::git+https://github.com/CavemanJay/gogurt.git")
sha256sums=("SKIP")

pkgver() {
  cd "$pkgname"

  # The most recent un-annotated tag
  # git describe --long --tags | sed 's/^v//;s/\([^-]*-g\)/r\1/;s/-/./g'

  # No tags
  printf "r%s.%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short HEAD)"
}

build() {
  export GOPATH="$srcdir"/gopath

  cd "$pkgname"

  go build -ldflags="-s -w" -o build/$pkgname
}

package() {
  cd "$pkgname"
  install -Dm755 ./build/gogurt "$pkgdir/usr/bin/gogurt"
}
