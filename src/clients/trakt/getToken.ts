import {
  traktApiUrl,
  traktApiVersion,
  traktRedictUri,
} from "@/utils/constants";
import { TraktAuthResponse } from "./models";

export const getTraktToken = async (
  clientId: string,
  clientSecret: string,
  code: string
): Promise<TraktAuthResponse> =>
  await fetch(traktApiUrl + "/oauth/token", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "trakt-api-version": traktApiVersion,
      "trakt-api-key": clientId,
    },
    body: JSON.stringify({
      client_id: clientId,
      client_secret: clientSecret,
      code: code,
      grant_type: "authorization_code",
      redirect_uri: traktRedictUri,
    }),
  })
    .then((res) => {
      if (!res.ok) {
        throw new Error(`Error fetching Trakt token: ${res.statusText}`);
      }
      return res;
    })
    .then((res) => res.json())
    .then((data) => data as TraktAuthResponse);
