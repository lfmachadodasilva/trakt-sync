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
import { Loader2Icon, Check, X } from "lucide-react";

export const Trakt = ({ cfg }: { cfg: ConfigEntity }) => {
  const clientIdRef = useRef<HTMLInputElement>(null);
  const clientSecretRef = useRef<HTMLInputElement>(null);
  const codeRef = useRef<HTMLInputElement>(null);
  const [saveStatus, setSaveStatus] = useState<
    "loading" | "success" | "error"
  >();

  const handleSave = () => {
    const updatedConfig: ConfigEntity = {
      trakt: {
        ...cfg.trakt,
        client_id: clientIdRef.current?.value || "",
        client_secret: clientSecretRef.current?.value || "",
        code: codeRef.current?.value || "",
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
  };

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Trakt Configuration</CardTitle>
        <CardDescription>Configure your Trakt settings</CardDescription>
      </CardHeader>
      <CardContent>
        <Label className="block mb-2 text-left">
          client ID
          <Input
            className="w-full mt-1 block w-full"
            type="text"
            defaultValue={cfg?.trakt?.client_id}
            autoComplete="off"
            ref={clientIdRef}
            placeholder="enter your client ID"
          />
        </Label>

        <Label className="block mb-2 text-left">
          client Secret
          <Input
            className="w-full mt-1 block w-full"
            // type="password"
            defaultValue={cfg?.trakt?.client_secret}
            autoComplete="off"
            ref={clientSecretRef}
            placeholder="enter your client secret"
          />
        </Label>
        <Label className="block mb-2 text-left">
          code
          <div className="flex items-center gap-2">
            <Input
              className="w-full mt-1 block w-full"
              // type="password"
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
        </Label>
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
