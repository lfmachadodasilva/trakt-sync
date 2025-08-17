import Database from 'better-sqlite3';
import path from 'path';
import fs from 'fs';
export * from "./config";

// Store the database instance with proper typing
let db: Database.Database | null = null;

function ensureConfigFolderExists(): void {
  const configPath = process.env.CONFIG_PATH || './config';
  if (!fs.existsSync(configPath)) {
    fs.mkdirSync(configPath, { recursive: true });
  }
}

function openDatabase(): Database.Database {
  ensureConfigFolderExists();
  const dbPath = path.join(process.env.CONFIG_PATH || './config', 'database.sqlite');
  return new Database(dbPath);
}

const isDatabaseOpen = (): boolean => {
  return db !== null && db.open;
};

export const getDatabase = () => {
  if (!isDatabaseOpen()) {
    db = openDatabase();

    // Initialize tables
    db.exec(
      "CREATE TABLE IF NOT EXISTS config (key TEXT PRIMARY KEY, value TEXT)"
    );
  }

  if (!db) {
    throw new Error("Database is not initialized");
  }

  return db;
};

export const closeDatabase = (): void => {
  if (db) {
    db.close();
    db = null;
  }
};


