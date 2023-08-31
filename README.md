# go-avito

<h3>Запуск проекта</h3>

Для запуска проекта необходимо запустить команду docker-compose up из корня проекта

<h2>Запросы</h2>

<h3>Пользователи</h3>

1. Получение пользователей

   GET http://localhost:8000/users
   Response
   Status code 200
   ```
   [
       {
           "id": 1,
           "name": "kabinochka",
           "email": "test@gmail.com"
       }
   ]
   ```

2. Получение пользователя по id

    GET http://localhost:8000/users/1  
    Response
    Status code 200
   ```
   {
      "id": 1,
      "name": "user_new",
      "email": "test@gmail.com"
   }
   ```

3. Создание пользователя

    POST http://localhost:8000/users

    Request

   ```
   {
       "name": "test2",
       "email": "test2@gmail.com"
   }
   ```

   Response

   Status code 204
   ```
   {
       "id": 3,
       "name": "test2",
       "email": "test2@gmail.com"
   }
   ```

4. Редактирование пользователя по id

   PUT http://localhost:8000/users/3

   Request

   ```
   {
       "name": "user_new",
       "email": "test@gmail.com"
   }
   ```

   Response

   204 No Content


5. Удаление пользователя по id
   DELETE http://localhost:8000/users/3

   Response
   ```
   "User deleted"
   ```

<h3>Сегменты</h3>

1. Получение сегментов
   GET http://localhost:8000/segments
   Response
   Status code 200
   ```
   [
      {
         "id": 1,
         "slug": "test_slug"
      }
   ]
   ```

2. Получение сегмента по id
   GET http://localhost:8000/segments/1
   Response
   Status code 200
   ```
   {
       "id": 1,
       "slug": "test_slug"
   }
   ```

3. Создание сегмента
   POST http://localhost:8000/users

   Request
   ```
   {
       "name": "test2",
       "email": "test2@gmail.com"
   }
   ```

   Response
   Status code 204
   ```
   {
       "id": 3,
       "name": "test2",
       "email": "test2@gmail.com"
   }
   ```

4. Редактирование сегмента по id
   PUT http://localhost:8000/segments/3

   Request
   ```
   {
       "slug": "test_slug_s"
   }
   ```

   Response
   ```
   204 No Content
   ```


5. Удаление сегмента по id
   DELETE http://localhost:8000/segments/3

   Response
   ```
   "Segment deleted"
   ```

<h3>Пользователи + сегменты</h3>
1. Добавление пользователя в сегмент
   POST http://localhost:8000/users/1/segments

   Request
   ```
   [
      {
         "slug": "test_slug_s"
      }
   ]
   ```

   Response
   200 OK
   ```
   "Segment added"
   [
      {
         "slug": "test_slug_s"
      }
   ]
   ```

2. Получение активных сегментов пользователя
   POST http://localhost:8000/users/1/segments

   Response
   200 OK
   ```
   [
       {
           "id": 1,
           "slug": "test_slug_s"
       }
   ]
   ```

3. Удаление пользователя из сегментов
   DELETE http://localhost:8000/users/1/segments

   Response
   200 OK
   ```
   "Segment deleted"
   [
       {
           "slug": "test_slug_s"
       }
   ]
   ```