FROM alpine:latest as build
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Europe/Prague /etc/localtime

FROM scratch as final
COPY --from=build /etc/localtime /etc/localtime
ADD /linux /
CMD ["/fischer_data_import_service"]