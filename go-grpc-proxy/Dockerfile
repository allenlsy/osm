FROM deis/go-dev

ARG XDS_TARGET
ENV GRPC_XDS_BOOTSTRAP "/bootstrap.json"
ENV GRPC_GO_LOG_VERBOSITY_LEVEL=99
ENV GRPC_GO_LOG_SEVERITY_LEVEL=info

COPY ./bin/xds-client-linux /xds-client
COPY ./sample-bootstrap.json /bootstrap.json


ENTRYPOINT [ "/xds-client" ]