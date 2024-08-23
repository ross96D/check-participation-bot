CREATE TABLE
  battle_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    position TEXT NOT NULL,
    fecha INTEGER NOT NULL,
    UNIQUE (position, fecha)
  );

CREATE TABLE
  player (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    team TEXT NOT NULL,
    UNIQUE (name, team)
  );

CREATE TABLE
  player_battle (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL,
    battle_log_id INTEGER NOT NULL,
    
    FOREIGN KEY (battle_log_id) REFERENCES battle_log (id),
    FOREIGN KEY (player_id) REFERENCES player (id)
  );

CREATE TABLE
  grupo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id INTEGER NOT NULL,
    
    UNIQUE(chat_id)
  );

CREATE TABLE
  grupo_battle (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    grupo_id INTEGER NOT NULL,
    battle_log_id INTEGER NOT NULL,
    FOREIGN KEY (battle_log_id) REFERENCES battle_log (id),
    FOREIGN KEY (grupo_id) REFERENCES grupo (id)
  );