import sqlite3 from "sqlite3";
import { Database, open } from "sqlite";
export * from "./config";

let db: Database<sqlite3.Database, sqlite3.Statement>;

async function openDb() {
  return open({
    filename: "./database.db",
    driver: sqlite3.Database,
  });
}

export const getDatabase = async () => {
  if (!db) {
    db = await openDb();

    await db.run(
      "CREATE TABLE IF NOT EXISTS config (key TEXT PRIMARY KEY, value TEXT)"
    );
  }
  return db;
};
