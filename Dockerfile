FROM scratch
ENV PGLET_SERVER_PORT=8080
EXPOSE 8080
COPY pglet /
ENTRYPOINT [ "./pglet", "server" ]