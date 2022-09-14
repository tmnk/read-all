CREATE TABLE items (
    id bigserial not null primary key,
    title varchar not null,
    description varchar not null,
    updated varchar DEFAULT null
);

INSERT INTO items (id, title, description, updated) VALUES
(1,	'database/sql',	'Рассказать про базы данных',	'rvasily'),
(2,	'memcache',	'Рассказать про мемкеш с примером использования',	NULL);

CREATE TABLE users (
  user_id bigserial not null primary key,
  login varchar not null,
  password varchar not null,
  email varchar not null,
  info varchar not null,
  updated varchar DEFAULT NULL
);

INSERT INTO users (user_id, login, password, email, info, updated) VALUES
(1,	'rvasily',	'love',	'rvasily@example.com',	'none',	NULL);