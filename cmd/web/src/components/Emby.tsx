import type { ConfigEntity, EmbyUser } from "@/config/models";
import { useEffect, useRef, useState } from "react";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "./ui/select";
import { getUsers } from "@/config/fetch";
import { Button } from "./ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";

export const Emby = ({ cfg }: { cfg: ConfigEntity }) => {
  const baseUrlRef = useRef<HTMLInputElement>(null);
  const apiKeyRef = useRef<HTMLInputElement>(null);
  const [users, setUsers] = useState<EmbyUser[]>([]);

  useEffect(() => {
    getUsers().then(setUsers);
  }, []);

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Emby Configuration</CardTitle>
        <CardDescription>Configure your Emby server settings</CardDescription>
      </CardHeader>
      <CardContent>
        <div key="base-url">
          <Label className="block mb-2">
            server url
            <Input
              type="text"
              defaultValue={cfg?.emby?.base_url}
              className="mt-1 block w-full p-2 border border-gray-300 rounded"
              placeholder="Enter your emby server URL"
              ref={baseUrlRef}
            />
          </Label>
        </div>
        <div key="api-key">
          <Label className="block mb-2">
            api key
            <Input
              type="text"
              defaultValue={cfg?.emby?.api_key}
              className="mt-1 block w-full p-2 border border-gray-300 rounded"
              placeholder="enter your emby api key"
              ref={apiKeyRef}
            />
          </Label>
        </div>
        <div key="user-id">
          <Label className="block mb-2">
            user id
            <Select defaultValue={cfg?.emby?.user_id || ""}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="select user" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Fruits</SelectLabel>

                  {users.map((user) => (
                    <SelectItem key={user.Id} value={user.Id}>
                      {user.Name}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </Label>
        </div>
      </CardContent>
      <CardFooter className="flex justify-end">
        <Button>save</Button>
      </CardFooter>
    </Card>
  );
};
