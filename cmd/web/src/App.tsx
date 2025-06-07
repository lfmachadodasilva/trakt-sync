import "./index.css";

import { Trakt } from "./components/Trakt";
import { Emby } from "./components/Emby";
import { useEffect, useState } from "react";
import { getConfig } from "./config/fetch";
import { Sync } from "./components/Sync";
import { AlertCircleIcon } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "./components/ui/alert";

export function App() {
  const [cfg, setCfg] = useState(null);
  const [error, setError] = useState<string | null>(null);
  const [refreshConfig, setRefreshConfig] = useState(false);

  useEffect(() => {
    getConfig()
      .then((config) => {
        setCfg(config);
      })
      .catch((error) => {
        console.error("Failed to fetch configuration:", error);
        setError("Failed to fetch configuration");
      });
  }, [refreshConfig]);

  const runRefresh = () => {
    setRefreshConfig((prev) => !prev);
  };

  return (
    <div className="container mx-auto p-8 text-center relative z-10 gap-8 flex flex-wrap justify-center items-start">
      {error && (
        <Alert variant="destructive" className="text-left w-full ">
          <AlertCircleIcon />
          <AlertTitle>Something went wrong</AlertTitle>
          <AlertDescription>
            Something went wrong while fetching the configuration. Please check
            your server logs for more details.
          </AlertDescription>
        </Alert>
      )}
      <Trakt cfg={cfg} refreshConfig={runRefresh} />
      <Emby cfg={cfg} refreshConfig={runRefresh} />
      <Sync cfg={cfg} refreshConfig={runRefresh} />
    </div>
    // <div className="container mx-auto p-8 text-center relative z-10">
    //   <div className="flex justify-center items-center gap-8 mb-8">
    //     <img
    //       src={logo}
    //       alt="Bun Logo"
    //       className="h-36 p-6 transition-all duration-300 hover:drop-shadow-[0_0_2em_#646cffaa] scale-120"
    //     />
    //     <img
    //       src={reactLogo}
    //       alt="React Logo"
    //       className="h-36 p-6 transition-all duration-300 hover:drop-shadow-[0_0_2em_#61dafbaa] [animation:spin_20s_linear_infinite]"
    //     />
    //   </div>
  );
}

export default App;
