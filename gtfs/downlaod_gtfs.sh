#!/bin/bash
wget -O bus.zip http://gtfs.ztp.krakow.pl/GTFS_KRK_A.zip
wget -O tram.zip http://gtfs.ztp.krakow.pl/GTFS_KRK_T.zip
unzip -u -d data/bus bus.zip
unzip -u -d data/tram tram.zip
