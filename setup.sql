START TRANSACTION;

CREATE TABLE players (
    id      VARCHAR(40) PRIMARY KEY
);

CREATE TABLE rooms (
    id      INT PRIMARY KEY,
    owner   VARCHAR(40),
    lives   INT,
    level   INT,
    FOREIGN KEY (owner) REFERENCES players(id)
);

CREATE TABLE in_game (
    roomId  INT,
    player  VARCHAR(40),
    PRIMARY KEY (player),
    FOREIGN KEY (roomId) REFERENCES rooms(id),
    FOREIGN KEY (player) REFERENCES players
);

CREATE TABLE cards(
    playerId    VARCHAR(20),
    card        INT,
    FOREIGN KEY (playerId) REFERENCES players(id)    
);