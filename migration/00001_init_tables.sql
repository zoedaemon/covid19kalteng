-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE "clients" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "name" varchar(255) NOT NULL,
    "key" varchar(255) NOT NULL,
    "secret" varchar(255) NOT NULL,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "roles" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "name" varchar(255) NOT NULL,
    "description" text,
    "system" varchar(255),
    "status" varchar(255),
    "permissions" varchar(255) ARRAY,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "users" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "roles" int ARRAY,
    "username" varchar(255) NOT NULL UNIQUE,
    "password" text NOT NULL,
    "email" varchar(255) UNIQUE,
    "phone" varchar(255) UNIQUE,
    "status" varchar(255),
    "first_login" boolean DEFAULT TRUE,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "edus" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "title" varchar(512),
    "description" text,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);


CREATE TABLE "cases" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "location" jsonb DEFAULT '{}',
    "data_detail" jsonb DEFAULT '[]',
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);


CREATE TABLE "data_definitions" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "key" varchar(255) NOT NULL,
    "name_total" varchar(255),
    "name_per_day" varchar(255),
    "stand_for" varchar(255) ,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);



-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS "clients" CASCADE;
DROP TABLE IF EXISTS "roles" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "edus" CASCADE;
DROP TABLE IF EXISTS "cases" CASCADE;