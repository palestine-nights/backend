CREATE DATABASE IF NOT EXISTS `restaurant`;
USE `restaurant`;

DROP TABLE IF EXISTS `reservations`;
DROP TABLE IF EXISTS `tables`;

CREATE TABLE IF NOT EXISTS `tables` (
  id          INT UNSIGNED NOT NULL AUTO_INCREMENT,
	places      TINYINT UNSIGNED NOT NULL,
	description VARCHAR(255) NOT NULL,
  created_at  DATETIME NOT NULL,
  updated_at  DATETIME NOT NULL,
  PRIMARY KEY (id)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `reservations` (
  id           INT UNSIGNED NOT NULL AUTO_INCREMENT,
  table_id     INT UNSIGNED NOT NULL,
  guests       TINYINT UNSIGNED NOT NULL,
  email        VARCHAR(63) NOT NULL,
  phone        VARCHAR(63) NOT NULL,
  status       ENUM('created', 'approved', 'cancelled'),
  fullname     VARCHAR(255) NOT NULL,
  time         DATETIME NOT NULL,
  created_at   DATETIME NOT NULL,
  updated_at   DATETIME NOT NULL,
  duration     BIGINT
  PRIMARY KEY (id),
  FOREIGN KEY (table_id)
    REFERENCES tables(id)
) ENGINE = InnoDB;
