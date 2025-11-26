-- Active: 1749016188917@@127.0.0.1@5432@agri_chain
-- PostgreSQL 数据库初始化脚本

-- CREATE DATABASE agri_chain;

-- 查看所有表
-- \dt

-- 删除表（如果存在）
DROP TABLE IF EXISTS consortium_nodes;
DROP TABLE IF EXISTS users;

-- 创建 consortium_nodes 表
CREATE TABLE consortium_nodes (
    id BIGINT PRIMARY KEY,
    addr VARCHAR(40),
    pub_key BYTEA,
    create_time BIGINT, 
    verify_time INTEGER
);

-- 创建 users 表
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(255),
    password VARCHAR(255)
);

-- 插入初始用户数据
INSERT INTO users(user_name, password) VALUES ('admin', 'agri_chain');
