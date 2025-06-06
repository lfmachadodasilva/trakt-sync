import type { ConfigEntity, EmbyUser } from "@/config/models";
import { use, useCallback, useEffect, useRef, useState } from "react";
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
import { getUsers, updateConfig } from "@/config/fetch";
import { Button } from "./ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Check, Loader2Icon, X } from "lucide-react";

export const Emby = ({ cfg }: { cfg: ConfigEntity }) => {
  const baseUrlRef = useRef<HTMLInputElement>(null);
  const apiKeyRef = useRef<HTMLInputElement>(null);
  const [users, setUsers] = useState<EmbyUser[]>([]);
  const [selectedUserId, setSelectedUserId] = useState<string>(
    cfg?.emby?.user_id || ""
  );
  const [saveStatus, setSaveStatus] = useState<
    "loading" | "success" | "error"
  >();

  useEffect(() => {
    if (
      (cfg?.emby?.base_url?.trim() ?? "") !== "" &&
      (cfg?.emby?.api_key?.trim() ?? "") !== ""
    ) {
      getUsers().then(setUsers);
    }
  }, [cfg?.emby?.base_url, cfg?.emby?.api_key]);

  useEffect(() => {
    if (cfg?.emby?.user_id) {
      setSelectedUserId(cfg.emby.user_id);
    }
  }, [cfg?.emby?.user_id]);

  const handleUserChange = (value: string) => setSelectedUserId(value);

  const handleSave = useCallback(() => {
    const updatedConfig: ConfigEntity = {
      emby: {
        ...cfg.emby,
        base_url: baseUrlRef.current?.value || "",
        api_key: apiKeyRef.current?.value || "",
        user_id: selectedUserId,
      },
    };

    setSaveStatus("loading");
    updateConfig(updatedConfig)
      .then(() => {
        console.debug("Configuration saved successfully");
        setSaveStatus("success");
      })
      .catch((error) => {
        console.debug("Failed to save configuration:", error);
        setSaveStatus("error");
      });
  }, [selectedUserId]);

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
            <Select value={selectedUserId} onValueChange={handleUserChange}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="select user" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Users</SelectLabel>

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
        <Button disabled={saveStatus === "loading"} onClick={handleSave}>
          {saveStatus === "loading" && <Loader2Icon className="animate-spin" />}
          save
          {saveStatus === "success" && <Check className="text-green-500" />}
          {saveStatus === "error" && <X className="text-red-500" />}
        </Button>
      </CardFooter>
    </Card>
  );
};
