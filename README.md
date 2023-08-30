# chat

1. POST users/add: Добавляет пользователя с именем user_name и возвращает его id. Если такой пользователь существует, возвращает ошибку.
```json
{
    "user_name": "irina"
}
```
``` json
{
    "id": 1
}
```

2. POST chats/add: Создает чат с заданным именем, состоящий из перечисленных пользователей и возращает его id. Если чат с таким именем существует или есть некоректные id пользователей, то возращает ошибку.
```json
{
    "name": "Photoshop",
    "users": ["1", "2"]
}
```
```json
{
    "id": 1
}
```

3. GET chats/get: Возвращает информацию о всех чатах пользователя. Если id пользователя некорректный, то возвращает ошибку. 
```json
{
    "id" 1
}
```
```json
{
    "chats": [
        {
            "id": 2,
            "name": "Blender",
            "users": [
                {
                    "id": 1,
                    "user_name": "irishka"
                }
            ]
        },
        {
            "id": 1,
            "name": "Photoshop",
            "users": [
                {
                    "id": 1,
                    "user_name": "irishka"
                },
                {
                    "id": 2,
                    "user_name": "ivan"
                }
            ]
        }
    ]
}
```

4. POST messages/add Добавляет сообщение в заданный чат от лица пользователя и возвращает его id. Если id пользователя или чата не существуют, то хендлер вернет ошибку.
```json
{
    "chat": "1", 
    "author": "2", 
    "text": "it's enough for today"
}
```
```json
{
    "id": 1
}
```

5. GET messages/get Возвращает список сообщений в заданном чате. Если id чата некорректный, то хендлер вернет ошибку.
```json
{
    "chat": 1
}
```
```json
{
    "messages": [
        {
            "id": 1,
            "chat": 1,
            "author": 2,
            "text": "teach me how to do this trick"
        },
        {
            "id": 2,
            "chat": 1,
            "author": 2,
            "text": "it's enough for today"
        },
        {
            "id": 3,
            "chat": 1,
            "author": 1,
            "text": "calm down. I'm super pro"
        }
    ]
}
```
    messages/get/debug Работает так же, только вместо id чата и автора сообщения возвращает их имена. 
```json
{
    "chat": "1"
}
```
```json
{
    "messages": [
        {
            "id": 1,
            "author": "ivan",
            "text": "teach me how to do this trick"
        },
        {
            "id": 2,
            "author": "ivan",
            "text": "it's enough for today"
        },
        {
            "id": 3,
            "author": "irishka",
            "text": "calm down. I'm super pro"
        }
    ]
}
```

6. POST users/add/tochat Добавляет пользователя в чат. Проверяет id пользователя и чата на корректность.
```json
{
    "user_id": "2",
    "chat_id": "2"
}
```
```json
{
    "message": "user has been added to chat"
}
```
