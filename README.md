# todo-list-v2-backend
I create todo list v2 to practice my web development skill.

Backend tech

Language : Go

Framework : Fiber

Database : MongoDB

Endpoints :
| Endpoint | Description | Required body |
| -------- | ----------- |---- |
| `[Post]` /api/register | register user | name, email, password |
| `[Post]` /api/login | login user | email, password |
| `[Get]` /api/user | get user's info | - |
| `[Post]` /api/logout | logout user | - |
| `[Post]` /api/password | change password | password |
| `[Post]` /api/task | add task | name, detail, time, priority, tags |
| `[Delete]` /api/task | delete task | _id (item id) |
