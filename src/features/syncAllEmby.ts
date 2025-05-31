import { getEmbyAllItems } from "@/clients/emby/getItems";
import { SyncData } from "./models";
import {
  TraktMarkAsWatchedRequest,
  TraktWatched,
  TraktWatchedResponse,
} from "@/clients/trakt/models";
import { EmbyItemResponse } from "@/clients/emby/models";
import { traktMoviesByImdbId } from "@/utils/trakt";
import { embyGetImdbId, embyItemsByImdbId } from "@/utils/emby";
import { markEmbyItemAsWatched } from "@/clients/emby/markWatched";
import { markTraktAsWatched } from "@/clients/trakt/markAsWatched";

const logPrefix = "sync all emby -";
export async function syncAllEmby(data: SyncData, trakt: TraktWatched) {
  if (!data.emby) {
    throw new Error("Emby data is required for syncAllEmby");
  }

  console.log(`${logPrefix} syncing all emby items with trakt...`);
  const emby = await getEmbyAllItems(
    data.emby.baseUrl,
    data.emby.apiKey,
    data.emby.userId
  );
  console.log(`${logPrefix} fetched emby items:`, {
    movies: emby.movies.length,
    series: emby.series.length,
  });
  await Promise.all([
    syncAllEmbyMovies(data, emby.movies, trakt.movies),
    syncAllEmbyShows(data, emby.series, trakt.shows),
  ]);
}

const syncAllEmbyMovies = async (
  data: SyncData,
  embyMovies: EmbyItemResponse[],
  traktMovies: TraktWatchedResponse[]
) => {
  if (!data.emby) {
    throw new Error("Emby data is required for syncAllEmbyMovies");
  }

  console.log(`${logPrefix} syncing emby movies with trakt...`);

  const trakt = traktMoviesByImdbId(traktMovies);
  const emby = embyItemsByImdbId(embyMovies);

  //   console.debug(`${logPrefix} emby movies: ${JSON.stringify(emby, null, 2)}`);

  const traktRequest: TraktMarkAsWatchedRequest = {};

  for (const key of Object.keys(emby)) {
    console.log(`${logPrefix} syncing movie with key: ${key}`);

    const traktMovie = trakt[key];
    const embyMovie = emby[key];

    if (!traktMovie && !embyMovie) {
      // If either does not exist, log the missing item
      console.error(`${logPrefix} missing both: ${key}`);
      continue;
    }

    if (!traktMovie && embyMovie.UserData.Played === true) {
      // Not watched on trakt, but watched in Emby
      // TODO - mark as watched in Trakt 🚧
      console.debug(
        `${logPrefix} trakt movie not found for emby movie ${embyMovie.Id} - ${embyMovie.Name}, but it is marked as watched in emby 🚧`
      );
      const imdbId = embyGetImdbId(embyMovie.ProviderIds);
      if (imdbId) {
        traktRequest.movies ??= [];
        traktRequest.movies.push({
          ids: { imdb: imdbId },
          watched_at: new Date(),
        });
      } else {
        console.warn(
          `${logPrefix} emby movie ${embyMovie.Id} - ${embyMovie.Name} does not have an IMDB ID, cannot mark as watched in Trakt.`
        );
      }

      continue;
    }

    if (traktMovie && embyMovie.UserData.Played === false) {
      // Watched on trakt, but not played in Emby

      try {
        const res = await markEmbyItemAsWatched(
          data.emby.baseUrl,
          data.emby.apiKey,
          data.emby.userId,
          embyMovie.Id
        );
        console.info(
          `${logPrefix} marked emby movie ${embyMovie.Id} - ${
            embyMovie.Name
          } as watched. Response: ${JSON.stringify(res, null, 2)}`
        );
      } catch (err) {
        console.error(
          `${logPrefix} error marking emby movie ${embyMovie.Id} - ${
            embyMovie.Name
          } as watched: ${err instanceof Error ? err.message : err}`
        );
      }
      continue;
    }

    console.debug(`${logPrefix} ${key} - ${embyMovie.Name} no action taken.`);
  }

  if (traktRequest.movies && traktRequest.movies.length > 0) {
    console.log(
      `${logPrefix} marking ${traktRequest.movies.length} movies as watched in Trakt...`
    );

    if (!data.trakt?.clientId || !data.trakt?.accessToken) {
      console.error(
        `${logPrefix} Trakt clientId or accessToken is missing, cannot mark movies as watched.`
      );
      return;
    }

    await markTraktAsWatched(
      traktRequest,
      data.trakt.clientId,
      data.trakt.accessToken
    );

    console.log(
      `${logPrefix} successfully marked ${
        traktRequest?.movies?.length ?? 0
      } movies as watched in Trakt.`
    );
  } else {
    console.log(`${logPrefix} no movies to mark as watched in Trakt.`);
  }
};

const syncAllEmbyShows = async (
  data: SyncData,
  embySeries: EmbyItemResponse[],
  traktShows: TraktWatchedResponse[]
) => {
  if (!data.emby) {
    throw new Error("Emby data is required for syncAllEmbyShows");
  }

  // TODO - Implement sync logic for Emby series with Trakt watched shows 🚧
  console.debug(
    `${logPrefix} syncing emby series with trakt ... ${data.emby.userId} ${embySeries.length} ${traktShows.length}`
  );
};
