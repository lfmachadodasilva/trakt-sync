import { getConfig, getDatabase, saveConfig } from "@/database";
import { Config } from "@/database/model";

export async function GET(request: Request) {
  //   await saveConfig({
  //     plex: {
  //       userId: "test plex",
  //     },
  //     trakt: {
  //       clientId: "test trakt",
  //     },
  //     emby: {
  //       userId: "test emby",
  //     },
  //   } as Config);
  console.log(getConfig());

  return Response.json({});
}
