import { EmbyConfig } from "@/clients/emby";
import { getDatabase } from ".";
import { Config } from "./model";
import { TraktConfig } from "@/clients/trakt";
import { PlexConfig } from "@/clients/plex";

export const getConfig = (): Config => {
  const db = getDatabase();

  const cfgRaw = db
    .prepare<{ key: string; value: string }[], { key: string; value: string }>(
      "SELECT * FROM config"
    )
    .all();

  const cfg: Config = {
    emby:
      (JSON.parse(
        cfgRaw.find((x) => x.key == "emby")?.value ?? ""
      ) as EmbyConfig) ?? {},
    trakt:
      (JSON.parse(
        cfgRaw.find((x) => x.key == "trakt")?.value ?? ""
      ) as TraktConfig) ?? {},
    plex:
      (JSON.parse(
        cfgRaw.find((x) => x.key == "plex")?.value ?? ""
      ) as PlexConfig) ?? {},
  };

  return cfg;
};

export const saveConfig = (cfg: Config): void => {
  const db = getDatabase();

  const insertQuery = `INSERT OR REPLACE INTO config (key, value) VALUES ${Object.entries(
    cfg
  )
    .map(() => "(?, ?)")
    .join(", ")}`;

  const values = Object.entries(cfg).flatMap(([key, value]) => [
    key,
    JSON.stringify(value) ?? {},
  ]);

  db.prepare(insertQuery).run(...values);
};
