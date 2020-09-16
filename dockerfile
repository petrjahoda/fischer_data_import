FROM alpine:latest as build
RUN apk add tzdata

FROM scratch as final
ADD /linux /
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
CMD ["/fischer_data_import_service_linux"]