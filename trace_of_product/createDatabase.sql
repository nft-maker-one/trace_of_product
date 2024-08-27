-- Active: 1715677941562@@127.0.0.1@3306@agri_chain
CREATE DATABASE IF NOT EXISTS agri_chain;
USE agri_chain;

DROP TABLE IF EXISTS consortium_nodes


SELECT * FROM consortium_nodes;
CREATE TABLE consortium_nodes (
    id BIGINT,
    addr VARCHAR(40),
    pub_key BLOB,
    create_time BIGINT
);

