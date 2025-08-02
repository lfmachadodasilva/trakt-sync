import { getDatabase } from "@/database";

export async function GET(request: Request) {
  const db = await getDatabase();
  return Response.json({});
}
