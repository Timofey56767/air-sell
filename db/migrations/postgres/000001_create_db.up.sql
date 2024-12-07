
CREATE TABLE IF NOT EXISTS "tickets" (
	"id" VARCHAR(255) NOT NULL UNIQUE, 
	"airline" varchar(255) NOT NULL,
	"departure_from" int NOT NULL,
	"arrival_at" int NOT NULL,
	"departure_time" timestamp NOT NULL,
	"arrival_time" timestamp NOT NULL,
	"quantity" int NOT NULL,
	"luggage" varchar(255) NOT NULL DEFAULT '',
	"hand_baggage" varchar(255) NOT NULL DEFAULT '',
	"price" decimal(10,2) NOT NULL,
	PRIMARY KEY("id")
);


CREATE TABLE IF NOT EXISTS "users" (
	"id" serial NOT NULL UNIQUE,
	"email" varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	PRIMARY KEY("id")
);


CREATE TABLE IF NOT EXISTS "buy_history" (
	"id" serial NOT NULL UNIQUE,
	"ticket_id" VARCHAR(255) NOT NULL,
	"user_id" int NOT NULL,
	"buy_time" timestamp NOT NULL,
	"count" int NOT NULL, 
	PRIMARY KEY("id")
);


CREATE TABLE IF NOT EXISTS "airports" (
	"id" serial NOT NULL UNIQUE,
	"city" varchar(255) NOT NULL,
	"iata_code" varchar(255) NOT NULL,
	"country" varchar(255) NOT NULL,
	PRIMARY KEY("id")
);


CREATE TABLE IF NOT EXISTS "passport_data" (
	"id" serial NOT NULL UNIQUE,
	"name" varchar(255) NOT NULL,
	"surname" varchar(255) NOT NULL,
	"patronymic" varchar(255) NOT NULL DEFAULT '',
	"passport_series_and_number" varchar(255) NOT NULL UNIQUE,
	"gender" char(1) NOT NULL CHECK(gender IN ('w', 'm')),
	"validity_period" date,
	"date_of_birth" date NOT NULL,
	"passport_type" varchar(255) NOT NULL CHECK(passport_type IN ('international_passport', 'passport')),
	"citizenship" varchar(255) NOT NULL,
	"user_id" int NOT NULL,
	PRIMARY KEY("id")
);


ALTER TABLE "buy_history"
ADD FOREIGN KEY("ticket_id") REFERENCES "tickets"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "buy_history"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "passport_data"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "tickets"
ADD FOREIGN KEY("departure_from") REFERENCES "airports"("id")
ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE "tickets"
ADD FOREIGN KEY("arrival_at") REFERENCES "airports"("id")
ON UPDATE RESTRICT ON DELETE RESTRICT;