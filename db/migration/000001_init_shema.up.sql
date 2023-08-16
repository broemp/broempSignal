CREATE TABLE "user" (
  "userid" BIGSERIAL PRIMARY KEY,
  "username" varchar NOT NULL,
  "discordid" bigint NOT NULL,
  "telegramid" bigint,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "group" (
  "groupid" BIGSERIAL PRIMARY KEY,
  "hostid" bigserial NOT NULL,
  "active" bool DEFAULT true,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "invite" (
  "groupid" bigserial,
  "guestid" bigserial,
  "accepted" bool DEFAULT null,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("groupid", "guestid")
);

CREATE TABLE "hoepperCount" (
  "hoepper" bigserial,
  "victim" bigserial,
  "hoepperCount" integer,
  PRIMARY KEY ("hoepper", "victim")
);

CREATE INDEX ON "user" ("discordid");

CREATE INDEX ON "group" ("groupid");

CREATE INDEX ON "group" ("hostid");

CREATE INDEX ON "invite" ("guestid");

CREATE INDEX ON "invite" ("groupid");

CREATE INDEX ON "hoepperCount" ("hoepper");

CREATE INDEX ON "hoepperCount" ("victim");

COMMENT ON COLUMN "invite"."accepted" IS 'null = waiting, false=declined, true=accepted';

ALTER TABLE "group" ADD FOREIGN KEY ("hostid") REFERENCES "user" ("userid");

ALTER TABLE "invite" ADD FOREIGN KEY ("groupid") REFERENCES "group" ("groupid");

ALTER TABLE "invite" ADD FOREIGN KEY ("guestid") REFERENCES "user" ("userid");

ALTER TABLE "hoepperCount" ADD FOREIGN KEY ("hoepper") REFERENCES "user" ("userid");

ALTER TABLE "hoepperCount" ADD FOREIGN KEY ("victim") REFERENCES "user" ("userid");
