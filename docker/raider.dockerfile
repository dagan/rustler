FROM dagan/rustler:latest AS rustler

FROM ubuntu:20.04
RUN apt-get update; apt-get install -y \
    curl=7.68.0-1ubuntu2.12 \
    dnsutils=1:9.16.1-0ubuntu2.10 \
    netcat=1.206-1ubuntu1 \
    net-tools=1.60+git20180626.aebd88e-1ubuntu1 \
    nmap=7.80+dfsg1-2build1 \
    && rm -rf /var/lib/apt/lists/*
RUN curl -L "https://dl.k8s.io/release/v1.22.12/bin/linux/amd64/kubectl" > /usr/bin/kubectl && chmod +x /usr/bin/kubectl
RUN useradd -ms /bin/bash -u 1001 raider
COPY --from=rustler /rustle /usr/bin/
COPY scripts/raider.sh /usr/bin/raider.sh
USER 1001
ENTRYPOINT ["/usr/bin/raider.sh"]