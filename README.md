# HTTP Запросы
### Получение баланса
`GET /users/{id}/balance`
### Начисление баланса
`PUT /users/{id}/balance/add`
### Списание баланса
`PUT /users/{id}/balance/reduce`
### Перевод баланса другому пользователю
`PUT /users/{sender_id}/balance/transfer/{receiver_id}?balance=`
