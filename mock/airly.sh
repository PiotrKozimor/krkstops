curl -X GET \
    --header 'Accept: application/json' \
    --header 'apikey: 8Q3j62QKaXmX0lT80DilDzvN6D2ESpOw' \
    'https://airapi.airly.eu/v2/measurements/installation?installationId=8077' > data/measurement.json

curl -X GET \
    --header 'Accept: application/json' \
    --header 'apikey: 8Q3j62QKaXmX0lT80DilDzvN6D2ESpOw' \
    'https://airapi.airly.eu/v2/installations/nearest?lat=50.062006&lng=19.940984&maxDistanceKM=5&maxResults=3' > data/nearest.json

curl -X GET \
    --header 'Accept: application/json' \
    --header 'apikey: 8Q3j62QKaXmX0lT80DilDzvN6D2ESpOw' \
    'https://airapi.airly.eu/v2/installations/8077' > data/installation.json