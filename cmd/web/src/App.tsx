import { Card, CardContent } from "@/components/ui/card";
import { APITester } from "./APITester";
import "./index.css";

import logo from "./logo.svg";
import reactLogo from "./react.svg";
import { Trakt } from "./components/Trakt";
import { Emby } from "./components/Emby";
import { useEffect, useState } from "react";
import { getConfig } from "./config/fetch";
import { Button } from "./components/ui/button";
import { Sync } from "./components/Sync";

export function App() {
  const [cfg, setCfg] = useState(null);

  useEffect(() => {
    getConfig().then((config) => {
      setCfg(config);
    });
  }, []);

  return (
    <div className="container mx-auto p-8 text-center relative z-10 gap-8 flex flex-wrap justify-center items-start">
      <Trakt cfg={cfg} />
      <Emby cfg={cfg} />
      <Sync cfg={cfg} />
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
