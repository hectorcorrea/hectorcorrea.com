CREATE DATABASE blogdb;

USE blogdb;

CREATE TABLE blogs (
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  summary VARCHAR(255) NULL,
  slug VARCHAR(255) NULL,
  content TEXT NULL,
  createdOn DATETIME NOT NULL,
  updatedOn DATETIME NULL,
  postedOn DATETIME NULL
);

CREATE INDEX blogs_index_title ON blogs(title);
CREATE INDEX blogs_index_postedOn ON blogs(postedOn DESC);


CREATE TABLE users (
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  login VARCHAR(255) NOT NULL,
  name VARCHAR(20) NOT NULL,
  password VARCHAR(255) NOT NULL
);

CREATE INDEX users_index_id ON users(id);


CREATE TABLE sessions (
  id char(64) NOT NULL PRIMARY KEY,
  userId INT NOT NULL,
  expiresOn DATETIME NOT NULL
);

CREATE INDEX sessions_index_id ON sessions(id);
