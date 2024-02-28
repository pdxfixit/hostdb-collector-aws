FROM scratch

LABEL maintainer="Ben Sandberg <info@pdxfixit.com>" \
      name="hostdb-collector-aws" \
      vendor="PDXfixIT, LLC"

COPY tls-ca-bundle.pem /etc/ssl/certs/

COPY hostdb-collector-aws /usr/bin/
COPY config.yaml /etc/hostdb-collector-aws/

ENTRYPOINT [ "/usr/bin/hostdb-collector-aws" ]
