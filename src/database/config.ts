import { EmbyConfig } from "@/clients/emby";
import { getDatabase } from ".";
import { Config } from "./model";
import { TraktConfig } from "@/clients/trakt";
import { PlexConfig } from "@/clients/plex";

export const getConfig = async (): Promise<Config> => {
  const db = await getDatabase();

  const cfgRaw = await db.all<{ key: string; value: string }[]>(
    "select * from config"
  );

  console.log(cfgRaw);

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

  console.log(cfg);

  return cfg;
};

export const saveConfig = async (cfg: Config) => {
  const db = await getDatabase();

  const stmt = await db.prepare(
    `INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)`
  );

  Object.entries(cfg).forEach(async ([key, value]) => {
    await stmt.run(key, JSON.stringify(value) ?? {});
  });
};
