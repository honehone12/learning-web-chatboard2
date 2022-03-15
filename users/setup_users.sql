DROP TABLE sessions;
DROP TABLE users;

CREATE TABLE users (
  id         SERIAL PRIMARY KEY,
  uu_id      VARCHAR(255) NOT NULL UNIQUE,
  name       VARCHAR(255),
  email      VARCHAR(255) NOT NULL UNIQUE,
  password   VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL   
);

CREATE TABLE sessions (
  id         SERIAL PRIMARY KEY,
  uu_id      VARCHAR(255) NOT NULL UNIQUE,
  user_name       VARCHAR(255),
  user_id    SERIAL REFERENCES users(id),
  created_at TIMESTAMP NOT NULL   
);
