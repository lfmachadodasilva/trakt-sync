import { getDatabase } from ".";
import { Config } from "./model";

interface ConfigRaw {
  key: string;
  value: string;
}

export const getConfig = async (): Promise<Config> => {
  const db = await getDatabase();

  const cfgRaw = await db.all<{ key: string; value: string }[]>(
    "select * from config"
  );
  return {} as Config;
};

export const saveConfig = async (config: Config) => {
  const db = await getDatabase();

  await db.run("INSERT INTO config (key, value) VALUES (?, ?)", []);
};
