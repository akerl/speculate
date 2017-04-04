FROM dock0/pkgforge
RUN pacman -S --needed --noconfirm go rsync
RUN curl -sLo /opt/x509.patch.b64 https://go-review.googlesource.com/changes/36941/revisions/1/patch?download && \
    base64 -d /opt/x509.patch.b64 > /opt/x509.patch && \
    rm /opt/x509.patch.b64
RUN curl -sLo /opt/go.tar.gz https://storage.googleapis.com/golang/go1.8.src.tar.gz && \
    tar -xvf /opt/go.tar.gz -C /opt && \
    cd /opt/go && \
    patch -p1 < /opt/x509.patch && \
    GOROOT_BOOTSTRAP=/usr/lib/go ./src/buildall.bash 'linux/amd64'
RUN ln -vs /opt/go/bin/{go,gofmt} /usr/local/bin/
    
