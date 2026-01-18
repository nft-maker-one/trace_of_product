-- Active: 1759935252193@@127.0.0.1@5432@agri_chain
-- PostgreSQL 数据库初始化脚本

-- CREATE DATABASE agri_chain;

-- 查看所有表
-- \dt

-- 删除表（如果存在）

CREATE DATABASE agri_chain;
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
-- 创建 users 表
CREATE TABLE users (
  "id" int4 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  "user_name" varchar(255),
  "nick_name" varchar(255),
  "email" varchar(255),
  "password" varchar(255),
  "invite_code" varchar(100),
  "agree_terms" bool DEFAULT false,
  "avatar_url" varchar(255) DEFAULT '/public/avatars/profile.jpg',
  "profile_public" bool DEFAULT true,
  "email_notifications" bool DEFAULT true,
  "last_login_at" int8,
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
);

-- 插入初始用户数据
INSERT INTO users(user_name, password) VALUES ('admin', 'agri_chain');
