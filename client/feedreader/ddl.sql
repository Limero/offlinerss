CREATE TABLE "CachedActions" (
  "action" INTEGER NOT NULL, "id" TEXT NOT NULL,
  "argument" INTEGER
);

CREATE TABLE "Enclosures" (
  "articleID" TEXT NOT NULL,
  "url" TEXT NOT NULL,
  "type" INTEGER NOT NULL,
  FOREIGN KEY(articleID) REFERENCES articles(articleID)
);

CREATE TABLE "articles" (
  "articleID" TEXT PRIMARY KEY NOT NULL UNIQUE,
  "feedID" TEXT NOT NULL, "title" TEXT NOT NULL,
  "author" TEXT, "url" TEXT NOT NULL,
  "html" TEXT NOT NULL, "preview" TEXT NOT NULL,
  "unread" INTEGER NOT NULL, "marked" INTEGER NOT NULL,
  "date" INTEGER NOT NULL, "guidHash" TEXT,
  "lastModified" INTEGER, "contentFetched" INTEGER NOT NULL
);

CREATE TABLE "categories" (
  "categorieID" TEXT PRIMARY KEY NOT NULL UNIQUE,
  "title" TEXT NOT NULL, "orderID" INTEGER,
  "exists" INTEGER, "Parent" TEXT, "Level" INTEGER
);

CREATE TABLE "feeds" (
  "feed_id" TEXT PRIMARY KEY NOT NULL UNIQUE,
  "name" TEXT NOT NULL, "url" TEXT NOT NULL,
  "category_id" TEXT, "subscribed" INTEGER DEFAULT 1,
  "xmlURL" TEXT, "iconURL" TEXT
);

CREATE TABLE "taggings" (
  "articleID" TEXT NOT NULL,
  "tagID" TEXT NOT NULL,
  FOREIGN KEY(articleID) REFERENCES articles(articleID),
  FOREIGN KEY(tagID) REFERENCES tags(tagID)
);

CREATE TABLE "tags" (
  "tagID" TEXT PRIMARY KEY NOT NULL UNIQUE,
  "title" TEXT NOT NULL, "exists" INTEGER,
  "color" INTEGER
);
