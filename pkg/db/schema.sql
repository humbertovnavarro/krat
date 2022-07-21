CREATE TABLE IF NOT EXISTS OnionServices (
    onionServiceID TEXT not null primary key,
    port INTEGER,
    privateKey BLOB
);

CREATE TABLE IF NOT EXISTS Logs (
    logLevel INTEGER not null,
    logEntry TEXT not null,
    createdAt INTEGER not null 
);