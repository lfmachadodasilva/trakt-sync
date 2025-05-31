import { EmbyWebhook } from "@/clients/emby/models";
import { markTraktAsWatched } from "@/clients/trakt/markAsWatched";
import { TraktMarkAsWatchedRequest } from "@/clients/trakt/models";
import { getConfig } from "@/data/config";
import { embyGetImdbId } from "@/utils/emby";

const logPrefix = "emby webhook - ";

export async function POST(request: Request) {
  const body: EmbyWebhook = await request.json();
  console.debug("Received webhook:", body);

  const finishToViewEvents =
    body.Event === "playback.stop" &&
    body.Title.includes("has finished playing");
  const markAsPlayedEvents = body.Event === "item.markplayed";

  if (!finishToViewEvents && !markAsPlayedEvents) {
    return new Response("Webhook received - no action taken");
  }

  // finish to view. mark as watched on trakt
  const imdbId = embyGetImdbId(body.Item?.ProviderIds);
  if (!imdbId) {
    console.error(`${logPrefix} no IMDB ID found for item:`, body.Item);
    return new Response("No IMDB ID found", { status: 400 });
  }

  switch (body.Item.Type) {
    case "Episode": {
      const episode = body.Item.IndexNumber ?? 0;
      const season = body.Item.ParentIndexNumber ?? 0;
      console.debug(`${logPrefix} Episode detected: S${season}E${episode}`);

      await processEpisode(imdbId, season, episode);
      break;
    }
    case "Movie": {
      console.debug(`${logPrefix} Movie detected:`, body.Item.Name);
      await processMovie(imdbId);
      break;
    }
    default: {
      console.warn(`${logPrefix} Unsupported item type:`, body.Item.Type);
      break;
    }
  }

  return new Response("Webhook received");
}

const processEpisode = async (
  imdbId: string,
  season: number,
  episode: number
) => {
  // Process episode logic here

  const dataSync = await getConfig();

  const request: TraktMarkAsWatchedRequest = {
    shows: [
      {
        ids: { imdb: imdbId },
        seasons: [
          {
            number: season,
            watched_at: new Date(),
            episodes: [
              {
                number: episode,
                watched_at: new Date(),
              },
            ],
          },
        ],
      },
    ],
  };

  if (!dataSync.trakt?.clientId || !dataSync.trakt?.accessToken) {
    console.error(`${logPrefix} Missing Trakt credentials`);
    throw new Error(`${logPrefix} Missing Trakt credentials`);
  }

  await markTraktAsWatched(
    request,
    dataSync.trakt.clientId,
    dataSync.trakt.accessToken
  );
};

const processMovie = async (imdbId: string) => {
  const dataSync = await getConfig();

  const request: TraktMarkAsWatchedRequest = {
    movies: [
      {
        ids: { imdb: imdbId },
        watched_at: new Date(),
      },
    ],
  };

  if (!dataSync.trakt?.clientId || !dataSync.trakt?.accessToken) {
    console.error(`${logPrefix} Missing Trakt credentials`);
    throw new Error(`${logPrefix} Missing Trakt credentials`);
  }

  await markTraktAsWatched(
    request,
    dataSync.trakt.clientId,
    dataSync.trakt.accessToken
  );
};
