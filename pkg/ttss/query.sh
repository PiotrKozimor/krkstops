stop=3338
t=$(date +%s )000
echo TRAM
curl "https://www.ttss.krakow.pl/internetservice/services/passageInfo/stopPassages/stop?stop=$stop&mode=departure&language=pl&timeFrame=40&startTime=$t"
echo BUS
curl "https://ttss.mpk.krakow.pl/internetservice/services/passageInfo/stopPassages/stop?stop=$stop&mode=departure&language=pl&timeFrame=40&startTime=$t"