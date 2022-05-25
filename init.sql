CREATE TABLE IF NOT EXISTS urls (
    short_url varchar(10) UNIQUE NOT NULL,
    long_url varchar UNIQUE NOT NULL,
    created_at timestamptz,
    CONSTRAINT PK_URL PRIMARY KEY (short_url, long_url)
);

CREATE INDEX shorturl_btree_IDX ON urls USING btree (short_url);
CREATE INDEX longurl_btree_IDX ON urls USING btree (long_url);