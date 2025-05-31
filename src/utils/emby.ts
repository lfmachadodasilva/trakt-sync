import { EmbyItemResponse } from "@/clients/emby/models";

export const embyGetImdbId = (
  providerIds: Record<string, string>
): string | null => {
  // Normalize keys to lowercase and check for "imdb"
  const imdbKey = Object.keys(providerIds).find(
    (key) => key.toLowerCase() === "imdb"
  );

  return imdbKey ? providerIds[imdbKey] : null;
};

export const embyItemsByImdbId = (
  embyMovies: EmbyItemResponse[]
): Record<string, EmbyItemResponse> =>
  embyMovies.reduce<Record<string, EmbyItemResponse>>((acc, movie) => {
    const imdbId = embyGetImdbId(movie.ProviderIds);
    if (imdbId) {
      acc[imdbId] = movie;
    }
    return acc;
  }, {});
