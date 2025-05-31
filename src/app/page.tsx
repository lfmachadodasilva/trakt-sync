"use client";

import { Emby } from "@/components/emby/emby";
import { Trakt } from "@/components/trakt/trakt";
import { SyncData } from "@/features/models";
import { syncAll } from "@/features/syncAll";

export default function Home() {
  const syncData = {
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
      redirectUrl: "urn:ietf:wg:oauth:2.0:oob",
    },
  } as SyncData;

  const handleSyncAll = async () => {
    await syncAll(syncData);
  };

  return (
    <div>
      <h1>Welcome to Trakt Sync</h1>

      <br></br>
      <Trakt data={syncData} />
      <br></br>
      <Emby data={syncData} />

      <br></br>
      <br></br>
      <button
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
        onClick={handleSyncAll}
      >
        Sync All
      </button>
    </div>
  );
}
