CREATE TABLE IF NOT EXISTS OnionServices (
    onionServiceID TEXT not null primary key,
    onionUrl TEXT not null,
    port INTEGER not null,
    privateKey TEXT not null
);

CREATE TABLE IF NOT EXISTS Logs (
    logLevel INTEGER not null,
    logEntry TEXT not null,
    createdAt INTEGER not null 
);