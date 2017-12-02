FROM dockerhub.cisco.com/spectre-docker-prod/spectre/alpinejre8:3.3_1.8.0
RUN apk add --no-cache bash
RUN apk add --no-cache wget
RUN apk add --update openssl
RUN mkdir -p /opt/docker/ssl
RUN /usr/bin/openssl req -x509 -newkey rsa:2048 -nodes -keyout /opt/docker/ssl/key.pem -out /opt/docker/ssl/cert.pem -days 365 -subj "/C=US/ST=California/L=San Jose/O=Cisco/OU=spectre/CN=be-txnaggregator"
ENV CERTFILE /opt/docker/ssl/cert.pem
ENV PKEYFILE /opt/docker/ssl/key.pem
ENV PROTOCOL https
COPY src/bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/.env /usr/bin/
COPY src/bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/scheduler-pagetest /usr/bin/
EXPOSE 4031
ENTRYPOINT ["/usr/bin/scheduler-pagetest", "-debug"]
