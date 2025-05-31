import { EmbyUserResponse } from "./models";

export const getEmbyUsers = async (
  baseUrl: string,
  apiKey: string
): Promise<EmbyUserResponse[]> => {
  return await fetch(`${baseUrl}/Users`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "X-Emby-Token": apiKey,
    },
  })
    .then((res) => {
      if (!res.ok) {
        throw new Error(`Failed to fetch users: ${res.statusText}`);
      }
      return res.json();
    })
    .then((data) => data as EmbyUserResponse[]);
};
