curl -X GET \
    'http://91.223.13.70/internetservice/services/passageInfo/stopPassages/stop?stop=81&mode=departure&language=pl' \
    > mock/bus1.json
curl -X GET \
    'http://91.223.13.70/internetservice/services/passageInfo/stopPassages/stop?stop=610&mode=departure&language=pl' \
    > mock/bus2.json
curl -X GET \
    'http://185.70.182.51/internetservice/services/passageInfo/stopPassages/stop?stop=610&mode=departure&language=pl' \
    > mock/tram2.json
# curl -X GET \
#     'http://91.223.13.70/internetservice/geoserviceDispatcher/services/stopinfo/stops?left=-648000000&bottom=-324000000&right=648000000&top=324000000' \
#     > mock/stops_bus.json
# curl -X GET \
#     'http://185.70.182.51/internetservice/geoserviceDispatcher/services/stopinfo/stops?left=-648000000&bottom=-324000000&right=648000000&top=324000000' \
#     > mock/stops_tram.json