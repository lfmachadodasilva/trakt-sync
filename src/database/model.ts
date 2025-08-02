import { EmbyConfig } from "@/clients/emby";
import { PlexConfig } from "@/clients/plex";
import { TraktConfig } from "@/clients/trakt";

export interface Config {
  trakt: TraktConfig;
  emby: EmbyConfig;
  plex: PlexConfig;
}
