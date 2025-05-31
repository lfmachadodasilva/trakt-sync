import {
  TraktMarkAsWatchedRequest,
  TraktMarkAsWatchedResponse,
} from "./models";

export const markTraktAsWatched = async (
  request: TraktMarkAsWatchedRequest,
  clientId: string,
  accessToken: string
): Promise<TraktMarkAsWatchedResponse> => {
  return await fetch("https://api.trakt.tv/sync/history", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "trakt-api-version": "2",
      "trakt-api-key": clientId,
      Authorization: `Bearer ${accessToken}`,
    },
    body: JSON.stringify(request),
  })
    .then((res) => {
      if (!res.ok) {
        throw new Error(`Failed to mark as watched: ${res.statusText}`);
      }
      return res.json();
    })
    .then((data) => data as TraktMarkAsWatchedResponse);
};
