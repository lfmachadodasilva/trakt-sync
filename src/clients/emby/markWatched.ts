import { EmbyPlayedItemsResponse } from "./models";

/**
 * Marks an Emby item as watched.
 * Documentation: https://swagger.emby.media/?staticview=true#/PlaystateService/postUsersByUseridPlayeditemsById
 *
 * @param baseUrl - The base URL of the Emby server.
 * @param token - The API token for authenticating with the Emby server.
 * @param userId - The ID of the user.
 * @param itemId - The ID of the item to mark as watched.
 * @returns A promise that resolves when the item has been marked as watched.
 */
export const markEmbyItemAsWatched = async (
  baseUrl: string,
  token: string,
  userId: string,
  itemId: string
): Promise<EmbyPlayedItemsResponse> =>
  await markEmbyItem(baseUrl, token, userId, itemId, true);

/**
 * Marks an Emby item as unwatched.
 * Documentation: https://swagger.emby.media/?staticview=true#/PlaystateService/deleteUsersByUseridPlayeditemsById
 *
 * @param baseUrl - The base URL of the Emby server.
 * @param token - The API token for authenticating with the Emby server.
 * @param userId - The ID of the user.
 * @param itemId - The ID of the item to mark as unwatched.
 * @returns A promise that resolves when the item has been marked as unwatched.
 */
export const markEmbyItemAsUnwatched = async (
  baseUrl: string,
  token: string,
  userId: string,
  itemId: string
): Promise<EmbyPlayedItemsResponse> =>
  await markEmbyItem(baseUrl, token, userId, itemId, false);

const markEmbyItem = async (
  baseUrl: string,
  token: string,
  userId: string,
  itemId: string,
  watched: boolean
): Promise<EmbyPlayedItemsResponse> => {
  if (!baseUrl || !token || !userId || !itemId) {
    throw new Error("All parameters are required.");
  }

  const url = `${baseUrl}/Users/${userId}/PlayedItems/${itemId}`;

  return await fetch(url, {
    method: watched ? "POST" : "DELETE",
    headers: {
      "Content-Type": "application/json",
      "X-Emby-Token": token,
    },
    body: JSON.stringify({}),
  })
    .then((res) => {
      if (!res.ok) {
        throw new Error(
          `Error marking emby item ${itemId} as ${
            watched ? "watched" : "unwatched"
          }: ${res.statusText}`
        );
      }
      return res.json();
    })

    .then((data) => data as EmbyPlayedItemsResponse);
};
