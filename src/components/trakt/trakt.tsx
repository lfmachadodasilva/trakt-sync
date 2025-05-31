// https://api.trakt.tv/oauth/authorize?response_type=code&client_id={{client_id}}&redirect_uri=https://display-parameters.com/

import { SyncData } from "@/features/models";

export const Trakt = ({ data }: { data: SyncData }) => {
  return (
    <div className="border-white border-2 p-4 rounded-lg">
      <h2>Trakt</h2>
      <br></br>

      <form className="mt-4">
        <label className="block mb-2">
          Client ID:
          <input
            type="text"
            value={data.trakt?.clientId}
            readOnly
            className="border border-gray-300 rounded p-2 w-full"
          />
        </label>
        <label className="block mb-2">
          Client Secrect:
          <input
            type="password"
            value={data.trakt?.clientSecret}
            className="border border-gray-300 rounded p-2 w-full"
          />
        </label>
        <label className="block mb-2">
          Code:
          <a
            href={`https://api.trakt.tv/oauth/authorize?response_type=code&client_id=${data.trakt?.clientId}&redirect_uri=${data.trakt?.redirectUrl}`}
            target="_blank"
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
          >
            get trakt code
          </a>
          <input
            type="password"
            value={data.trakt?.code}
            className="border border-gray-300 rounded p-2 w-full"
          />
        </label>
        <label className="block mb-2">
          Access token:
          <input
            type="password"
            value={data.trakt?.accessToken}
            className="border border-gray-300 rounded p-2 w-full"
          />
        </label>
        <label className="block mb-2">
          Refresh token:
          <input
            type="password"
            value={data.trakt?.refreshToken}
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
