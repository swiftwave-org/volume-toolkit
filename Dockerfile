FROM busybox

RUN mkdir /app
RUN chmod 777 /app
RUN mkdir /backup
RUN chmod 777 /backup

COPY ./volume_toolkit /usr/bin/volume_toolkit
RUN chmod +x /usr/bin/volume_toolkit

ENTRYPOINT ["/usr/bin/volume_toolkit"]