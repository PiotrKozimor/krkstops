set -xe
for t in A M T
do 
    f=GTFS_KRK_$t.zip
    curl https://gtfs.ztp.krakow.pl/$f > $f
    unzip -o $f -d GTFS_KRK_$t
    curl https://gtfs.ztp.krakow.pl/TripUpdates_$t.pb > TripUpdates_$t.pb        
done
curl https://gtfs.ztp.krakow.pl/TripUpdates.pb > TripUpdates.pb