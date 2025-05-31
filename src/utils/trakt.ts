import { TraktWatchedResponse } from "@/clients/trakt/models";

export const traktMoviesByImdbId = (
  traktMovies: TraktWatchedResponse[]
): Record<string, TraktWatchedResponse> =>
  traktMovies.reduce<Record<string, TraktWatchedResponse>>((acc, movie) => {
    const imdbId = movie.movie?.ids?.imdb;
    if (imdbId) {
      acc[imdbId] = movie;
    }
    return acc;
  }, {});
