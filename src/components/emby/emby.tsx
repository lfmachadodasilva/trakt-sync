import { getEmbyUsers } from "@/clients/emby/getUsers";
import { EmbyUserResponse } from "@/clients/emby/models";
import { SyncData } from "@/features/models";
import { onSubmitAction } from "./action";

export const Emby = async ({ data }: { data: SyncData }) => {
  console.log("Emby users1:", data.emby?.userId);
  const users: EmbyUserResponse[] = data.emby?.baseUrl
    ? await getEmbyUsers(data.emby.baseUrl, data.emby.apiKey)
    : [];

  if (data.emby && users.length === 0) {
    data.emby.userId = "";
  } else if (data.emby && !data.emby.userId && users.length > 0) {
    data.emby.userId = users[0].Id;
  }

  console.log("Emby users2:", data.emby?.userId);

  return (
    <div className="border-white border-2 p-4 rounded-lg">
      <h2>Emby</h2>
      <br></br>

      <form className="mt-4" action={onSubmitAction}>
        <label className="block mb-2">
          Base URL:
          <input
            name="baseUrl"
            type="text"
            defaultValue={data.emby?.baseUrl}
            className="border border-gray-300 rounded p-2 w-full"
          />
        </label>
        <label className="block mb-2">
          User ID:
          <select
            name="userId"
            defaultValue={data.emby?.userId}
            className="border border-gray-300 rounded p-2 w-full"
          >
            {users.map((user) => (
              <option key={user.Id} value={user.Id}>
                {user.Name}
              </option>
            ))}
          </select>
        </label>
        <label className="block mb-2">
          API Key:
          <input
            name="apiKey"
            type="password"
            defaultValue={data.emby?.apiKey}
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
