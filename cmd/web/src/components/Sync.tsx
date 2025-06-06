import type { ConfigEntity } from "@/config/models";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Button } from "./ui/button";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { useEffect, useMemo, useRef, useState } from "react";
import { runSync, updateConfig } from "@/config/fetch";
import { Check, Loader2Icon, X } from "lucide-react";
import cronstrue from "cronstrue";

export const Sync = ({ cfg }: { cfg: ConfigEntity }) => {
  const cronRef = useRef<HTMLInputElement>(null);
  const [saveStatus, setSaveStatus] = useState<
    "loading" | "success" | "error"
  >();
  const [syncStatus, setSyncStatus] = useState<
    "loading" | "success" | "error"
  >();
  const [cronText, setCronText] = useState<string>(null);

  useEffect(() => {
    setCronText(getCronjobText(cfg?.cronjob));
    cronRef.current.value = cfg?.cronjob || "";
  }, [cfg?.cronjob]);

  const getCronjobText = (cron: string): string => {
    if (!cron || cron.trim() === "") {
      return "No cron expression provided";
    }
    try {
      return cronstrue.toString(cron);
    } catch (error) {
      console.debug("Invalid cron expression:", error);
      return "Invalid cron expression";
    }
  };

  const handleCronChange = () => {
    const cronValue = cronRef.current?.value || "";
    setCronText(getCronjobText(cronValue));
  };

  const handleSave = () => {
    const cronValue = cronRef.current?.value || "";

    const updatedConfig: ConfigEntity = {
      cronjob: cronValue,
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
  };

  const handleSync = () => {
    setSyncStatus("loading");
    runSync()
      .then(() => {
        console.debug("Sync started successfully");
        setSyncStatus("success");
      })
      .catch((error) => {
        console.debug("Failed to start sync:", error);
        setSyncStatus("error");
      });
  };

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Sync</CardTitle>
        <CardDescription>Configure your sync settings</CardDescription>
        <CardAction>
          <Button
            variant="secondary"
            disabled={syncStatus === "loading"}
            onClick={handleSync}
          >
            {syncStatus === "loading" && (
              <Loader2Icon className="animate-spin" />
            )}
            run async
            {syncStatus === "success" && <Check className="text-green-500" />}
            {syncStatus === "error" && <X className="text-red-500" />}
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>
        <div key="base-url">
          <Label className="block mb-2">
            cron job
            <Input
              readOnly
              type="text"
              defaultValue={cfg?.cronjob}
              className="mt-1 block w-full p-2 border border-gray-300 rounded"
              placeholder="Enter your cron job frequency (e.g., '0 * * * *')"
              ref={cronRef}
              onChange={handleCronChange}
            />
            <Label className="text-sm text-muted-foreground mt-1">
              This will run: <strong>{cronText}</strong>
            </Label>
            <br></br>
            <br></br>
            <Label className="text-sm text-muted-foreground mt-1">
              This cron job will run based on your defined schedule. Need help?
              Visit{" "}
              <a
                href="https://crontab.guru"
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-500 underline"
              >
                crontab.guru
              </a>
            </Label>
          </Label>
        </div>
        {/* <pre>{JSON.stringify(cfg, null, 2)}</pre> */}
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
