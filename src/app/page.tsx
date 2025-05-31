import { getEmbyAllItems } from "@/clients/emby/getItems";
import { getTraktAllWatched } from "@/clients/trakt/getWatched";
import { Emby } from "@/components/emby/emby";
import { RunAll } from "@/components/runAll";
import { Trakt } from "@/components/trakt/trakt";
import { getConfig, upsertConfig } from "@/data/config";
import { SyncData } from "@/features/models";
import { embyGetImdbId } from "@/utils/emby";
import { stringify } from "@/utils/stringify";

export default async function Home() {
  let syncData: SyncData;

  try {
    syncData = await getConfig();
  } catch {
    syncData = {
      emby: {
        userId: "aac3a78d9f184ea480fb1629e76aad57",
        apiKey: "b039ba2b065e4ba1bca2307cce593478",
        baseUrl: "http://192.168.1.13:8096",
      },
      trakt: {
        clientId:
          "eb4ede9a384157e9aa60aad8c72c36c0485215659c82ad7b1fe965359a55caf4",
        clientSecret:
          "0b2df529b1b229102030549ec0d76480f36f50cfcc94e695a6f6bd43994a6d17",
        accessToken:
          "fb386d3c6fcbf20104a33b0687953e43ade6469dff4123bdd032eb88f7d53d1c",
        refreshToken:
          "b185099900ad6bd02f331dd35b0a6d72e0ff599384a46b5b5e894ed64d9f353d",
        redirectUrl: "urn:ietf:wg:oauth:2.0:oob",
      },
    } as SyncData;
    try {
      await upsertConfig(syncData);
    } catch {
      console.error(
        "Failed to save default sync data. Please check your database connection."
      );
    }
  }

  const emby =
    syncData.emby &&
    (await getEmbyAllItems(
      syncData.emby.baseUrl,
      syncData.emby.apiKey,
      syncData.emby.userId
    ));
  const trakt =
    syncData.trakt &&
    (await getTraktAllWatched(
      syncData.trakt.accessToken,
      syncData.trakt.clientId
    ));

  return (
    <div>
      <h1>Welcome to Trakt Sync</h1>

      <br></br>
      <Trakt data={syncData} />
      <br></br>
      <Emby data={syncData} />
      <br></br>
      <RunAll data={syncData} />

      <div className="border-white border-2 p-4 rounded-lg">
        <h2>debug</h2>
        {emby && (
          <textarea
            defaultValue={stringify(
              emby.series.filter(
                (x) => embyGetImdbId(x.ProviderIds) === "tt0386676"
              )
            )}
            readOnly
            style={{ width: "50%", height: "300px" }}
          />
        )}
        {trakt && (
          <textarea
            defaultValue={stringify(
              trakt.shows.filter((x) => x.show.ids.imdb.startsWith("tt0386676"))
            )}
            readOnly
            style={{ width: "50%", height: "300px" }}
          />
        )}
      </div>
    </div>
  );
}
