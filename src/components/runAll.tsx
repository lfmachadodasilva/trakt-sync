"use client";

import { SyncData } from "@/features/models";
import { syncAll } from "@/features/syncAll";

export const RunAll = ({ data }: { data: SyncData }) => {
  const handleSyncAll = async () => {
    await syncAll(data);
  };

  return (
    <button
      className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
      onClick={handleSyncAll}
    >
      Sync All
    </button>
  );
};
