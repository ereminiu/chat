create table users (
    id serial primary key, 
    name varchar(255)
);

create table chats (
    id serial primary key,
    name varchar(255)
);

create table users_to_chats (
    id serial primary key, 
    user_id int not null references users,
    chat_id int not null references chats,
    unique (user_id, chat_id)
);

create table messages (
    id serial primary key, 
    msg text 
);

create table messages_to_chats (
    id serial primary key, 
    message_id int not null references messages,
    chat_id int not null references chats,
    unique (message_id, chat_id)
);

create table messages_to_users (
    id serial primary key, 
    message_id int not null references messages, 
    user_id int not null references users,
    unique (message_id, user_id)
);