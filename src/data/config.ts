import { SyncData } from "@/features/models";
import { getDbClient } from "./client";

export const upsertConfig = async (data: SyncData) => {
  const client = await getDbClient();
  await client.query(`
    CREATE TABLE IF NOT EXISTS config
    (
        id character varying(100) NOT NULL,
        data json,
        PRIMARY KEY (id)
    );
  `);

  const upsertQuery = async (id: string, data: object) =>
    await client.query(
      "INSERT INTO config (id, data) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data;",
      [id, data]
    );

  const promises = [];
  if (data.trakt) {
    promises.push(upsertQuery("trakt", data.trakt));
  }
  if (data.emby) {
    promises.push(upsertQuery("emby", data.emby));
  }
  if (data.plex) {
    promises.push(upsertQuery("plex", data.plex));
  }
  if (data.jellyfin) {
    promises.push(upsertQuery("jellyfin", data.jellyfin));
  }
  await Promise.all(promises);
};

export const getConfig = async (): Promise<SyncData> => {
  const client = await getDbClient();
  const res = await client.query("SELECT * FROM config");
  const config: SyncData = {};

  res.rows.forEach((row) => {
    config[row.id as keyof SyncData] = row.data;
  });

  return config;
};
