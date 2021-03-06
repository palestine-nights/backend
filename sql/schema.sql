CREATE DATABASE IF NOT EXISTS `restaurant`;
USE `restaurant`;

DROP TABLE IF EXISTS `reservations`;
DROP TABLE IF EXISTS `tables`;
DROP TABLE IF EXISTS `menu`;
DROP TABLE IF EXISTS `categories`;

CREATE TABLE IF NOT EXISTS `tables` (
  id          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  places      TINYINT UNSIGNED NOT NULL,
  description VARCHAR(255) NOT NULL,
  active      BOOLEAN NOT NULL DEFAULT TRUE,
  updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `reservations` (
  id           INT UNSIGNED NOT NULL AUTO_INCREMENT,
  table_id     INT UNSIGNED NOT NULL,
  guests       TINYINT UNSIGNED NOT NULL,
  email        VARCHAR(63) NOT NULL,
  phone        VARCHAR(63) NOT NULL,
  state        ENUM('created', 'approved', 'cancelled') NOT NULL DEFAULT 'created',
  full_name    VARCHAR(255) NOT NULL,
  time         DATETIME NOT NULL,
  duration     BIGINT,
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (table_id)
    REFERENCES tables(id)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `categories` (
  id          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name        VARCHAR(255) NOT NULL,
  `order`     INT UNSIGNED NOT NULL,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `menu` (
  id          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name        VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  price       FLOAT NOT NULL,
  category_id INT UNSIGNED NOT NULL,
  image_url   TEXT NOT NULL,
  active      BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (category_id)
    REFERENCES categories(id)
) ENGINE = InnoDB;

