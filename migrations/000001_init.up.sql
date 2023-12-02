CREATE TABLE IF NOT EXISTS "user"(
    "id" serial NOT NULL PRIMARY KEY,
    "username" varchar(30)  UNIQUE NOT NULL,
    "password" text NOT NULL,
    "first_name" varchar(15) NOT NULL,
    "last_name" varchar(15) NOT NULL,
    "phone_number" varchar(15),
    "mail" text
);

CREATE TABLE IF NOT EXISTS "static"(
    "id" serial NOT NULL PRIMARY KEY,
    "user_id" int,
    "battery_charge" float,
    "battery_status" float,
    "data_trans_stand" float,
    "sim_presence" float
);

CREATE TABLE IF NOT EXISTS "dynamic"(
    "id" serial NOT NULL PRIMARY KEY,
    "user_id" int,
    "max_device_offs" float,
    "min_device_offs" float,
    "max_dev_acceleration" float,
    "min_dev_acceleration" float,
    "min_light" float,
    "max_light" float,
    "hit_y" float,
    "hit_x" float
);