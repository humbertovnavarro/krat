CREATE TABLE IF NOT EXISTS OnionServices (
    onionServiceID TEXT not null primary key,
    port INTEGER,
    privateKey BLOB
)
