"use server";

import { SyncData } from "@/features/models";
import { syncAll } from "@/features/syncAll";

export const handleSyncAllAction = async (syncData: SyncData) => {
  syncAll(syncData)
    .then(() => {
      console.log("Sync completed successfully!");
    })
    .catch((error) => {
      console.error("Error during sync:", error);
    });
};
