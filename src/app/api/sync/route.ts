import { getConfig, saveConfig } from "@/database";
import { Config } from "@/database/model";

export async function GET(request: Request) {
  
  const cfg = await getConfig();

  return Response.json(cfg);
}

export async function POST(request: Request) {
  await saveConfig({
    plex: {
      userId: "test plex",
    },
    trakt: {
      clientId: "test trakt",
    },
    emby: {
      userId: "test emby",
    },
  } as Config);
  return Response.json({});
}
