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


CREATE TABLE sessions (
  id char(64) NOT NULL PRIMARY KEY,
  expiresOn DATETIME NOT NULL
);

CREATE INDEX sessions_index_id ON sessions(id);
