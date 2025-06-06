import type { ConfigEntity } from "@/config/models";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { Button } from "./ui/button";
import { useRef, useState } from "react";
import { updateConfig } from "@/config/fetch";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";

export const Trakt = ({ cfg }: { cfg: ConfigEntity }) => {
  const clientIdRef = useRef<HTMLInputElement>(null);
  const clientSecretRef = useRef<HTMLInputElement>(null);
  const codeRef = useRef<HTMLInputElement>(null);

  const handleSave = async () => {
    const updatedConfig: ConfigEntity = {
      ...cfg,
      trakt: {
        ...cfg.trakt,
        client_id: clientIdRef.current?.value || "",
        client_secret: clientSecretRef.current?.value || "",
        code: codeRef.current?.value || "",
      },
    };

    try {
      await updateConfig(updatedConfig);
    } catch (error) {
      console.error("Failed to save configuration:", error);
    }
  };

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Trakt Configuration</CardTitle>
        <CardDescription>Configure your Trakt settings</CardDescription>
      </CardHeader>
      <CardContent>
        <div>
          <Label>client ID</Label>
          <Input
            className="w-full"
            type="text"
            defaultValue={cfg?.trakt?.client_id}
            autoComplete="off"
            ref={clientIdRef}
            placeholder="enter your client ID"
          />
        </div>
        <div>
          <Label>client Secret</Label>
          <Input
            className="w-full"
            type="password"
            defaultValue={cfg?.trakt?.client_secret}
            autoComplete="off"
            ref={clientSecretRef}
            placeholder="enter your client secret"
          />
        </div>
        <div>
          <Label>code</Label>
          <div className="flex items-center gap-2">
            <Input
              className="w-full"
              type="password"
              defaultValue={cfg?.trakt?.code}
              ref={codeRef}
              placeholder="enter your code"
            />
            <Button asChild className="btn btn-primary">
              <a href={cfg?.trakt?.redirect_url} target="_blank">
                get code
              </a>
            </Button>
          </div>
        </div>
      </CardContent>
      <CardFooter className="flex justify-end">
        <Button onClick={handleSave}>save</Button>
      </CardFooter>
    </Card>
  );
};
