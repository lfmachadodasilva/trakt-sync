import type { ConfigEntity, EmbyUser } from "./models";

const API_URL =
  process.env.NODE_ENV === "development" ? "http://localhost:4000" : "";

export const getConfig = async (): Promise<ConfigEntity> => {
  try {
    const response = await fetch(`${API_URL}/config`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error("Failed to fetch config:", error);
    throw error;
  }
};

export const updateConfig = async (config: ConfigEntity): Promise<void> => {
  try {
    const response = await fetch(`${API_URL}/config`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(config),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
  } catch (error) {
    console.error("Failed to update config:", error);
    throw error;
  }
};

export const getUsers = async (): Promise<EmbyUser[]> => {
  try {
    const response = await fetch(`${API_URL}/emby/users`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error("Failed to fetch users:", error);
    throw error;
  }
};
