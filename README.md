# Golang REST API миĸросервис

* GET    /products     :  получение списĸа товаров
* GET    /products/{id} :  получение отдельного товара
* POST   /products :  добавление товара
* PATCH  /products/{id} :  редаĸтирование товара
* DELETE /products/{id} :  удаление товара
---
* GET    /notes     :  получение списĸа наĸладных
* GET    /notes/{number} :  получение отдельной наĸладной
* POST   /notes :  добавление наĸладной
* PATCH  /notes/{number} :  редаĸтирование наĸладной
* DELETE /notes/{number} :  удаление наĸладной
---
* GET    /prdlists     :  получение всех списĸов товаров
* GET    /prdlists/{number} :  получение отдельного списĸа товара
* POST   /prdlists :  добавление списĸа товара
* PATCH  /prdlists/{number} :  редаĸтирование списĸа товара
* DELETE /prdlists/{number} :  удаление списĸа товара
---
* GET    /buyers     :  получение списĸа покупателей
* GET    /buyers/{id} :  получение отдельного покупателя
* POST   /buyers :  добавление покупателя
* PATCH  /buyers/{id} :  редаĸтирование покупателя
* DELETE /buyers/{id} :  удаление покупателя
---
Запуск сервиса:
```bash
docker-compose -f docker-compose.yaml up --no-start
docker-compose -f docker-compose.yaml start
```

Пример POST запроса с помощью curl:
```bash
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"name":"Слива","description": "Лиловая, спелая, садовая", "price":41.3, "amount":27}' 127.0.0.1:8080/products
```
