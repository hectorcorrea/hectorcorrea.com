CREATE DATABASE hectorcorrea;

USE hectorcorrea;

CREATE TABLE blogs (
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  summary VARCHAR(255) NULL,
  url VARCHAR(255) NULL,
  createdOn DATETIME NOT NULL,
  updatedOn DATETIME NULL,
  postedOn DATETIME NULL
);

CREATE INDEX blogs_index_title ON blogs(title);
CREATE INDEX blogs_index_postedOn ON blogs(postedOn DESC);
