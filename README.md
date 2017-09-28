## InstalLation
```
docker build -t fx-service .
docker run -it -p 8000:8000 -e "APP_ID={YOUR_OPEN_EXCHANGE_APP_ID}" --name running-fx-service fx-service
```
