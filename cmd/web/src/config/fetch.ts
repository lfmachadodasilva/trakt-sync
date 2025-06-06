import type { ConfigEntity, EmbyUser } from "./models";

const API_URL =
  process.env.NODE_ENV === "development" ? "http://localhost:4000" : "";

export const getConfig = async (): Promise<ConfigEntity> =>
  new Promise<ConfigEntity>(async (resolve, reject) => {
    try {
      const response = await fetch(`${API_URL}/config`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) {
        reject(new Error(`HTTP error! status: ${response.status}`));
      }
      resolve(await response.json());
    } catch (error) {
      console.error("Failed to fetch config:", error);
      reject(error);
    }
  });

export const updateConfig = async (config: ConfigEntity): Promise<void> =>
  new Promise<void>(async (resolve, reject) => {
    try {
      const response = await fetch(`${API_URL}/config`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(config),
      });
      if (!response.ok) {
        reject(new Error(`HTTP error! status: ${response.status}`));
      }
      resolve();
    } catch (error) {
      console.error("Failed to update config:", error);
      reject(error);
    }
  });

export const getUsers = async (): Promise<EmbyUser[]> =>
  new Promise<EmbyUser[]>(async (resolve, reject) => {
    try {
      const response = await fetch(`${API_URL}/emby/users`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) {
        reject(new Error(`HTTP error! status: ${response.status}`));
      }
      resolve(await response.json());
    } catch (error) {
      console.error("Failed to fetch users:", error);
      reject(error);
    }
  });

export const runSync = async (): Promise<void> =>
  new Promise<void>(async (resolve, reject) => {
    try {
      const response = await fetch(`${API_URL}/sync`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) {
        reject();
        // throw new Error(`HTTP error! status: ${response.status}`);
      }
      resolve();
    } catch (error) {
      //   console.error("Failed to run sync:", error);
      //   throw error;
      reject();
    }
  });
