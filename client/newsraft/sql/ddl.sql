CREATE TABLE "feeds" (
  feed_url TEXT NOT NULL UNIQUE,
  title TEXT,
  link TEXT,
  content TEXT,
  attachments TEXT,
  persons TEXT,
  extras TEXT,
  download_date INTEGER NOT NULL DEFAULT 0,
  update_date INTEGER NOT NULL DEFAULT 0,
  time_to_live INTEGER NOT NULL DEFAULT 0,
  http_header_etag TEXT,
  http_header_last_modified INTEGER NOT NULL DEFAULT 0,
  http_header_expires INTEGER NOT NULL DEFAULT 0,
  user_data TEXT
);

CREATE TABLE "items" (
  feed_url TEXT NOT NULL,
  guid TEXT NOT NULL UNIQUE, -- unique manually added
  title TEXT,
  link TEXT,
  content TEXT,
  attachments TEXT,
  persons TEXT,
  extras TEXT,
  publication_date INTEGER NOT NULL DEFAULT 0,
  update_date INTEGER NOT NULL DEFAULT 0,
  unread INTEGER NOT NULL DEFAULT 0,
  important INTEGER NOT NULL DEFAULT 0,
  user_data TEXT,
  download_date INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_items_eight_way ON items(
  feed_url, guid, title, link, publication_date,
  update_date, unread, important
);

CREATE INDEX idx_items_guid ON items(feed_url, guid);
