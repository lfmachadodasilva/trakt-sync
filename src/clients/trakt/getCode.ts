// https://api.trakt.tv/oauth/authorize?response_type=code&client_id=eb4ede9a384157e9aa60aad8c72c36c0485215659c82ad7b1fe965359a55caf4&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob

import { traktApiUrl, traktRedictUri } from "@/utils/constants";

export const getTraktCodeUrl = async (clientId: string): Promise<string> => {
  const url = traktApiUrl + "/oauth/authorize";
  const redirectUri = traktRedictUri;
  const responseType = "code";

  return `${url}?response_type=${responseType}&client_id=${clientId}&redirect_uri=${encodeURIComponent(
    redirectUri
  )}`;
};
