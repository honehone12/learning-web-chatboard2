DROP TABLE posts;
DROP TABLE threads;
DROP TABLE sessions;
DROP TABLE users;

CREATE TABLE users (
  id         SERIAL PRIMARY KEY,
  uu_id      VARCHAR(255) NOT NULL UNIQUE,
  name       VARCHAR(255) NOT NULL UNIQUE,
  email      VARCHAR(255) NOT NULL UNIQUE,
  password   VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL   
);

CREATE TABLE sessions (
  id         SERIAL PRIMARY KEY,
  uu_id      VARCHAR(255) NOT NULL UNIQUE,
  user_name  VARCHAR(255),
  user_id    SERIAL REFERENCES users(id),
  state      TEXT,
  created_at TIMESTAMP NOT NULL   
);

CREATE TABLE threads (
  id          SERIAL PRIMARY KEY,
  uu_id       VARCHAR(255) NOT NULL UNIQUE,
  topic       TEXT,
  num_replies SERIAL,
  owner       VARCHAR(255),
  user_id     SERIAL REFERENCES users(id),
  last_update TIMESTAMP NOT NULL,
  created_at  TIMESTAMP NOT NULL       
);

CREATE TABLE posts (
  id          SERIAL PRIMARY KEY,
  uu_id       VARCHAR(255) NOT NULL UNIQUE,
  body        TEXT,
  contributor VARCHAR(255),
  user_id     SERIAL REFERENCES users(id),
  thread_id   SERIAL REFERENCES threads(id),
  created_at  TIMESTAMP NOT NULL  
);
