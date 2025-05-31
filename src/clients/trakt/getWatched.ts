import { traktApiUrl, traktApiVersion } from "@/utils/constants";
import { TraktWatched, TraktWatchedResponse } from "./models";

/**
 * Fetches the list of all watched movies and shows from Trakt API.
 * Documentation: https://trakt.docs.apiary.io/#reference/sync/get-watched/get-watched
 *
 * @param token - The API token for authenticating with the Trakt API.
 * @param clientId - The client ID for the Trakt API.
 * @returns A promise that resolves to an object containing arrays of watched movies and shows.
 */
export const getTraktAllWatched = async (
  token: string,
  clientId: string
): Promise<TraktWatched> => {
  const [movies, shows] = await Promise.all([
    getTraktMoviesWatched(token, clientId),
    getTraktShowsWatched(token, clientId),
  ]);

  return { movies, shows } as TraktWatched;
};

const getTraktMoviesWatched = async (
  token: string,
  clientId: string
): Promise<TraktWatchedResponse[]> =>
  await getTraktWatched("movies", token, clientId);

const getTraktShowsWatched = async (
  token: string,
  clientId: string
): Promise<TraktWatchedResponse[]> =>
  await getTraktWatched("shows", token, clientId);

const getTraktWatched = async (
  type: "movies" | "shows",
  token: string,
  clientId: string
): Promise<TraktWatchedResponse[]> => {
  return await fetch(`${traktApiUrl}/sync/watched/${type}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      "trakt-api-version": traktApiVersion,
      "trakt-api-key": clientId,
    },
    cache: "force-cache",
    next: {
      // Revalidate every 60 seconds
      revalidate: 60,
      tags: [`trakt:${type}:watched`],
    },
  })
    .then((res) => {
      if (!res.ok) {
        throw new Error(
          `Error fetching Trakt watched ${type}: ${res.statusText}`
        );
      }
      return res.json();
    })
    .then((data) => data as TraktWatchedResponse[]);
};
