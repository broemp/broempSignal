CREATE TABLE "user" (
  "discordid" bigint PRIMARY KEY,
  "username" varchar NOT NULL,
  "telegramid" bigint,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "group" (
  "groupid" BIGSERIAL PRIMARY KEY,
  "hostid" bigint NOT NULL,
  "active" bool DEFAULT true,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "invite" (
  "groupid" bigserial,
  "guestid" bigint,
  "accepted" bool DEFAULT null,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("groupid", "guestid")
);

CREATE TABLE "hoepperCount" (
  "hoepper" bigint,
  "victim" bigint,
  "hoepperCount" integer,
  PRIMARY KEY ("hoepper", "victim")
);

CREATE TABLE "afk" (
  "afkid" BIGSERIAL PRIMARY KEY,
  "userid" bigint,
  "created_at" timestamp DEFAULT (now())
);

CREATE INDEX ON "user" ("discordid");

CREATE INDEX ON "group" ("groupid");

CREATE INDEX ON "group" ("hostid");

CREATE INDEX ON "invite" ("guestid");

CREATE INDEX ON "invite" ("groupid");

CREATE INDEX ON "hoepperCount" ("hoepper");

CREATE INDEX ON "hoepperCount" ("victim");

CREATE INDEX ON "afk" ("userid");

COMMENT ON COLUMN "invite"."accepted" IS 'null = waiting, false=declined, true=accepted';

ALTER TABLE "group" ADD FOREIGN KEY ("hostid") REFERENCES "user" ("discordid");

ALTER TABLE "invite" ADD FOREIGN KEY ("groupid") REFERENCES "group" ("groupid");

ALTER TABLE "invite" ADD FOREIGN KEY ("guestid") REFERENCES "user" ("discordid");

ALTER TABLE "hoepperCount" ADD FOREIGN KEY ("hoepper") REFERENCES "user" ("discordid");

ALTER TABLE "hoepperCount" ADD FOREIGN KEY ("victim") REFERENCES "user" ("discordid");

ALTER TABLE "afk" ADD FOREIGN KEY ("userid") REFERENCES "user" ("discordid");
