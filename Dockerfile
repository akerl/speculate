FROM dock0/pkgforge
RUN pacman -S --needed --noconfirm go rsync
ENV PATCH_URL https://go-review.googlesource.com/changes/36941/revisions/1/patch?download
ENV GO_VERSION 1.8.1
ENV GO_URL https://storage.googleapis.com/golang/go${GO_VERSION}.src.tar.gz
RUN curl -sLo /opt/x509.patch.b64 $PATCH_URL && \
    base64 -d /opt/x509.patch.b64 > /opt/x509.patch && \
    rm /opt/x509.patch.b64 && \
    curl -sLo /opt/go.tar.gz $GO_URL && \
    tar -xvf /opt/go.tar.gz -C /opt && \
    cd /opt/go && \
    patch -p1 < /opt/x509.patch && \
    cd src && \
    GOROOT_BOOTSTRAP=/usr/lib/go ./buildall.bash 'linux/amd64' && \
    ln -vs /opt/go/bin/{go,gofmt} /usr/local/bin/
