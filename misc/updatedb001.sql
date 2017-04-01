USE blogdb;

ALTER TABLE blogs ADD COLUMN contentMd TEXT NULL;

UPDATE blogs SET contentMd = content;
