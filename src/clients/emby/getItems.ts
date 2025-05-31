import { EmbyItemResponse, EmbyItems } from "./models";

/**
 * Fetches all items (movies and series) from the Emby server.
 * Open api: https://swagger.emby.media/?staticview=true
 * Api keys: https://dev.emby.media/doc/restapi/API-Key-Authentication.html
 *
 * @param baseUrl - The base URL of the Emby server (e.g., "http://localhost:8096")
 * @param token - The API token for authenticating with the Emby server
 * @param userId - The ID of the user whose items are being fetched
 * @returns A promise that resolves to an object containing arrays of movies and series
 */
export const getEmbyAllItems = async (
  baseUrl: string,
  token: string,
  userId: string
): Promise<EmbyItems> => {
  const [movies, series] = await Promise.all([
    getEmbyMovies(baseUrl, token, userId),
    getEmbySeries(baseUrl, token, userId),
  ]);

  await Promise.all(
    series.map(async (s) => {
      s.Episodes = await getEmbyEpisodes(baseUrl, token, userId, s.Id);
    })
  );

  return {
    movies,
    series,
  } as EmbyItems;
};

export const getEmbyMovies = async (
  baseUrl: string,
  token: string,
  userId: string
): Promise<EmbyItemResponse[]> => getEmbyItems("Movie", baseUrl, token, userId);

const getEmbySeries = async (
  baseUrl: string,
  token: string,
  userId: string
): Promise<EmbyItemResponse[]> =>
  getEmbyItems("Series", baseUrl, token, userId);

const getEmbyItems = async (
  type: "Movie" | "Series",
  baseUrl: string,
  token: string,
  embyUserId: string
): Promise<EmbyItemResponse[]> => {
  if (!baseUrl || !token) {
    throw new Error("Base URL and token are required to fetch Emby items.");
  }

  // Remove trailing slash from baseUrl if present
  const sanitizedBaseUrl = baseUrl.endsWith("/")
    ? baseUrl.slice(0, -1)
    : baseUrl;

  // Build query parameters
  const queryParams = new URLSearchParams({
    IncludeItemTypes: type,
    Recursive: "true",
    Fields: "ProviderIds",
  });

  const url = `${sanitizedBaseUrl}/Users/${embyUserId}/Items?${queryParams.toString()}`;

  return fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "X-Emby-Token": token,
    },
    // cache: "no-cache", // Ensure fresh data
  })
    .then((res) => {
      if (!res.ok) {
        throw new Error(`Error fetching Emby ${type}: ${res.statusText}`);
      }
      return res.json();
    })
    .then((data) => data.Items as EmbyItemResponse[]);
};

const getEmbyEpisodes = async (
  baseUrl: string,
  token: string,
  embyUserId: string,
  seriesId: string
): Promise<EmbyItemResponse[]> => {
  if (!baseUrl || !token) {
    throw new Error("Base URL and token are required to fetch Emby episodes.");
  }
  // Remove trailing slash from baseUrl if present
  const sanitizedBaseUrl = baseUrl.endsWith("/")
    ? baseUrl.slice(0, -1)
    : baseUrl;
  // Build query parameters
  const queryParams = new URLSearchParams({
    Recursive: "true",
    Fields: "ProviderIds",
    EnableUserData: "true",
    UserId: embyUserId,
  });

  const url = `${sanitizedBaseUrl}/Shows/${seriesId}/Episodes?${queryParams.toString()}`;

  return fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "X-Emby-Token": token,
    },
  })
    .then((res) => {
      if (!res.ok) {
        throw new Error(`Error fetching Emby episodes: ${res.statusText}`);
      }
      return res.json();
    })
    .then((data) => data.Items as EmbyItemResponse[]);
};
