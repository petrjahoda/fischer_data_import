#!/usr/bin/env bash
cd linux || exit
upx fischer_data_import_service_linux
cd ..
docker rmi -f petrjahoda/fischer_data_import_service:latest
docker build -t petrjahoda/fischer_data_import_service:latest .
docker push petrjahoda/fischer_data_import_service:latest

docker rmi -f petrjahoda/fischer_data_import_service:2020.3.3
docker build -t petrjahoda/fischer_data_import_service:2020.3.3 .
docker push petrjahoda/fischer_data_import_service:2020.3.3
