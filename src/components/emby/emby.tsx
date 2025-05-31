import { SyncData } from "@/features/models";

export const Emby = ({ data }: { data: SyncData }) => {
  return (
    <div className="border-white border-2 p-4 rounded-lg">
      <h2>Emby</h2>
      <br></br>

      <form className="mt-4">
        <label className="block mb-2">
          Base URL:
          <input
            type="text"
            value={data.emby?.baseUrl}
            readOnly
            className="border border-gray-300 rounded p-2 w-full"
          />
        </label>
        <label className="block mb-2">
          User ID:
          <input
            type="text"
            value={data.emby?.userId}
            readOnly
            className="border border-gray-300 rounded p-2 w-full"
          />
        </label>
        <label className="block mb-2">
          API Key:
          <input
            type="password"
            value={data.emby?.apiKey}
            readOnly
            className="border border-gray-300 rounded p-2 w-full"
          />
        </label>
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
        >
          Save
        </button>
      </form>
    </div>
  );
};
