import type { ConfigEntity, EmbyUser } from "@/config/models";
import { useCallback, useEffect, useRef, useState } from "react";
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
import {
  Check,
  ListRestart,
  Loader2Icon,
  RefreshCcw,
  Save,
  X,
} from "lucide-react";

export const Emby = ({
  cfg,
  refreshConfig,
}: {
  cfg: ConfigEntity;
  refreshConfig: () => void;
}) => {
  const baseUrlRef = useRef<HTMLInputElement>(null);
  const apiKeyRef = useRef<HTMLInputElement>(null);
  const [users, setUsers] = useState<EmbyUser[]>([]);
  const [refetchUsers, setRefreshUsers] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState<string>(
    cfg?.emby?.user_id || ""
  );
  const [saveStatus, setSaveStatus] = useState<
    "loading" | "success" | "error"
  >();
  const [resetStatus, setResetStatus] = useState<
    "loading" | "success" | "error"
  >();

  useEffect(() => {
    if (
      (cfg?.emby?.base_url?.trim() ?? "") !== "" &&
      (cfg?.emby?.api_key?.trim() ?? "") !== ""
    ) {
      getUsers()
        .then(setUsers)
        .catch((error) => {
          console.debug("Failed to fetch users:", error);
          setUsers([]);
        });
    }
  }, [cfg?.emby?.base_url, cfg?.emby?.api_key]);

  useEffect(() => {
    if (cfg?.emby?.user_id) {
      setSelectedUserId(cfg.emby.user_id);
    }
  }, [cfg?.emby?.user_id, refetchUsers]);
  useEffect(() => {
    if (users.length > 0 && (!selectedUserId || selectedUserId === "")) {
      setSelectedUserId(users[0].Id);
    }
  }, [users]);

  const handleUserChange = (value: string) => setSelectedUserId(value);

  const handleRefetchUsers = useCallback(() => {
    handleSave().then(() => {
      setRefreshUsers((prev) => !prev);
    });
  }, []);

  const handleReset = () => {
    setResetStatus("loading");
    updateConfig({ ...cfg, emby: {} } as ConfigEntity)
      .then(() => {
        console.debug("Configuration saved successfully");
        setResetStatus("success");
        refreshConfig();
      })
      .catch((error) => {
        console.error("Failed to save configuration:", error);
        setResetStatus("error");
      });
  };

  const handleSave = useCallback(async () => {
    const updatedConfig: ConfigEntity = {
      emby: {
        ...(cfg?.emby ?? {}),
        base_url: baseUrlRef.current?.value || "",
        api_key: apiKeyRef.current?.value || "",
        user_id: selectedUserId,
      },
    };

    setSaveStatus("loading");
    return await updateConfig(updatedConfig)
      .then(() => {
        setSaveStatus("success");
        refreshConfig();
      })
      .catch((error) => {
        console.error("Failed to save configuration:", error);
        setSaveStatus("error");
      });
  }, [selectedUserId, cfg]);

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>emby configuration</CardTitle>
        <CardDescription>configure your emby server settings</CardDescription>
      </CardHeader>
      <CardContent>
        <div key="base-url">
          <Label className="block mb-2 text-left">
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
          <Label className="block mb-2 text-left">
            api key
            <Input
              type="text"
              defaultValue={cfg?.emby?.api_key}
              className="mt-1 block w-full p-2 border border-gray-300 rounded"
              placeholder="enter your emby api key"
              ref={apiKeyRef}
            />
            {cfg?.emby?.base_url && (
              <CardDescription className="text-left">
                You can find your API key in the{" "}
                <a
                  href={`${cfg.emby.base_url}/web/index.html#!/apikeys`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="underline"
                >
                  advanced settings
                </a>{" "}
                of your emby server.
              </CardDescription>
            )}
          </Label>
        </div>
        <div key="user-id">
          <Label className="block mb-2 text-left">
            user
            <div className="flex items-center gap-2">
              <Select value={selectedUserId} onValueChange={handleUserChange}>
                <SelectTrigger className="w-full mt-1">
                  <SelectValue placeholder="select user" className="" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectLabel>users</SelectLabel>
                    {users.map((user) => (
                      <SelectItem key={user.Id} value={user.Id}>
                        {user.Name}
                      </SelectItem>
                    ))}
                  </SelectGroup>
                </SelectContent>
              </Select>
              <Button className="btn btn-primary" onClick={handleRefetchUsers}>
                <RefreshCcw />
              </Button>
            </div>
          </Label>
        </div>
      </CardContent>
      <CardFooter className="flex justify-end gap-4">
        {/* <Button variant="outline">
          <TestTubeDiagonal /> test
        </Button> */}
        <Button
          variant="outline"
          disabled={resetStatus === "loading"}
          onClick={handleReset}
        >
          {resetStatus === "loading" && (
            <Loader2Icon className="animate-spin" />
          )}
          <ListRestart />
          reset
          {resetStatus === "success" && <Check className="text-green-500" />}
          {resetStatus === "error" && <X className="text-red-500" />}
        </Button>
        <Button disabled={saveStatus === "loading"} onClick={handleSave}>
          {saveStatus === "loading" && <Loader2Icon className="animate-spin" />}
          <Save />
          save
          {saveStatus === "success" && <Check className="text-green-500" />}
          {saveStatus === "error" && <X className="text-red-500" />}
        </Button>
      </CardFooter>
    </Card>
  );
};
