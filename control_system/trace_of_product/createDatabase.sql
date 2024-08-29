-- Active: 1715677941562@@127.0.0.1@3306@agri_chain
CREATE DATABASE IF NOT EXISTS agri_chain;
USE agri_chain;
SHOW TABLES;

DROP TABLE IF EXISTS consortium_nodes;

SELECT * FROM consortium_nodes;
CREATE TABLE consortium_nodes (
    id BIGINT,
    addr VARCHAR(40),
    pub_key BLOB,
    create_time BIGINT
);
DROP TABLE IF EXISTS users;
SELECT * FROM users;
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_name VARCHAR(255),
    password VARCHAR(255)
);


