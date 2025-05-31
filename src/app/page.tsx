"use client";

import { SyncData } from "@/features/models";
import { syncAll } from "@/features/syncAll";

export default function Home() {
  const embyUserId = "aac3a78d9f184ea480fb1629e76aad57";
  const embyApiKey = "b039ba2b065e4ba1bca2307cce593478";
  const embyBaseUrl = "http://192.168.1.13:8096";
  const traktClientId =
    "eb4ede9a384157e9aa60aad8c72c36c0485215659c82ad7b1fe965359a55caf4";
  const traktToken =
    "fb386d3c6fcbf20104a33b0687953e43ade6469dff4123bdd032eb88f7d53d1c";
  const traktRedirectUrl = "urn:ietf:wg:oauth:2.0:oob";

  const syncData = {
    emby: {
      userId: embyUserId,
      apiKey: embyApiKey,
      baseUrl: embyBaseUrl,
    },
    trakt: {
      clientId: traktClientId,
      token: traktToken,
      redirectUrl: traktRedirectUrl,
    },
  } as SyncData;

  const handleSyncAll = async () => {
    await syncAll(syncData);
  };

  return (
    <div>
      <h1>Welcome to Trakt Sync</h1>

      <button
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
        onClick={handleSyncAll}
      >
        Sync All
      </button>
    </div>
  );
}
