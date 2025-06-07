import type { ConfigEntity } from "@/config/models";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { Button } from "./ui/button";
import { useEffect, useRef, useState } from "react";
import { getTraktCodeUrl, setTraktCode, updateConfig } from "@/config/fetch";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";
import {
  Loader2Icon,
  Check,
  X,
  SquareArrowOutUpRight,
  Save,
  ListRestart,
} from "lucide-react";

export const Trakt = ({
  cfg,
  refreshConfig,
}: {
  cfg: ConfigEntity;
  refreshConfig: () => void;
}) => {
  const clientIdRef = useRef<HTMLInputElement>(null);
  const clientSecretRef = useRef<HTMLInputElement>(null);
  const codeRef = useRef<HTMLInputElement>(null);
  const [saveStatus, setSaveStatus] = useState<
    "loading" | "success" | "error"
  >();
  const [resetStatus, setResetStatus] = useState<
    "loading" | "success" | "error"
  >();
  const [codeUrl, setCodeUrl] = useState<string>();

  useEffect(() => {
    getTraktCodeUrl()
      .then((url) => {
        setCodeUrl(url);
      })
      .catch((error) => {
        console.debug("Failed to fetch Trakt code URL:", error);
        setCodeUrl("");
      });
  });

  const handleReset = () => {
    setResetStatus("loading");
    updateConfig({
      ...cfg,
      trakt: {
        ...cfg.trakt,
        access_token: null,
        refresh_token: null,
        code: null,
      },
    } as ConfigEntity)
      .then(() => {
        console.debug("Configuration saved successfully");
        setResetStatus("success");
        refreshConfig();
      })
      .catch((error) => {
        console.debug("Failed to save configuration:", error);
        setResetStatus("error");
      });
  };

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
        if (
          updatedConfig.trakt.code &&
          updatedConfig.trakt.code !== "" &&
          (!updatedConfig.trakt.access_token ||
            updatedConfig.trakt.access_token === "")
        ) {
          setTraktCode(updatedConfig.trakt.code)
            .then(() => {
              console.debug("Trakt code set successfully");
              setSaveStatus("success");
            })
            .catch((error) => {
              console.debug("Failed to set Trakt code:", error);
              setSaveStatus("error");
            });
        } else {
          setSaveStatus("success");
        }
      })
      .catch((error) => {
        console.debug("Failed to save configuration:", error);
        setSaveStatus("error");
      });
  };

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>trakt Configuration</CardTitle>
        <CardDescription>configure your trakt settings</CardDescription>
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
              <a href={codeUrl} target="_blank">
                <SquareArrowOutUpRight />
              </a>
            </Button>
          </div>
        </Label>
        <CardDescription className="text-left">
          create a new application on trakt in{" "}
          <a
            href="https://trakt.tv/oauth/applications"
            target="_blank"
            className="underline"
          >
            here
          </a>
          . remember to set the redirect URL to{" "}
          <code>urn:ietf:wg:oauth:2.0:oob</code>
        </CardDescription>
      </CardContent>
      <CardFooter className="flex justify-end gap-4">
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
