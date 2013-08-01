CREATE TABLE item (is_new NUMERIC, desc BLOB, id INTEGER PRIMARY KEY, link TEXT, rss_id NUMERIC, title TEXT, dir TEXT);
CREATE TABLE rss (updated_at TEXT, lang TEXT, id INTEGER PRIMARY KEY, desc TEXT, title TEXT, link TEXT, url TEXT);
