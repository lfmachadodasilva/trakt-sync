"use server";

import { getTraktAllWatched } from "@/clients/trakt/getWatched";
import { SyncData } from "./models";
import { syncAllEmby } from "./syncAllEmby";

export async function syncAll(data: SyncData) {
  if (!data.trakt) {
    throw new Error("Trakt data is required for syncAll");
  }

  console.log("starting sync all...");
  const trakt = await getTraktAllWatched(data.trakt.token, data.trakt.clientId);
  console.log("fetched trakt watched items:", {
    movies: trakt.movies.length,
    shows: trakt.shows.length,
  });

  if (data.emby) {
    await syncAllEmby(data, trakt);
  }
  if (data.plex) {
    // TODO - Implement Plex sync 🚧
  }
  if (data.jellyfin) {
    // TODO - Implement Jellyfin sync 🚧
  }
}
