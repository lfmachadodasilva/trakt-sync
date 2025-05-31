"use server";

import { upsertConfig } from "@/data/config";
import { SyncData } from "@/features/models";
import { revalidatePath } from "next/cache";

export const onSubmitAction = async (formData: FormData) => {
  const data = Object.fromEntries(formData);

  await upsertConfig({
    emby: {
      userId: data.userId as string,
      apiKey: data.apiKey as string,
      baseUrl: data.baseUrl as string,
    },
  } as SyncData);

  revalidatePath("/");
};
