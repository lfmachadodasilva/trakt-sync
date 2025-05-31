import { Client } from "pg";

let client: Client | null = null;

export const getDbClient = async (): Promise<Client> => {
  if (!client) {
    client = new Client(
      "postgres://user:password@localhost:5432/trakt-sync?sslmode=disable"
    );
    await client.connect();
  }

  return client;
};
