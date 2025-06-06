export interface TraktConfig {
  client_id?: string;
  client_secret?: string;
  access_token?: string;
  refresh_token?: string;
  code?: string;
  redirect_url?: string;
  expires_in?: number;
  created_at?: number;
}

export interface EmbyConfig {
  user_id?: string;
  api_key?: string;
  base_url?: string;
}

export interface PlexConfig {
  user_id?: string;
}

export interface JellyfinConfig {
  user_id?: string;
}

export interface ConfigEntity {
  trakt?: TraktConfig;
  emby?: EmbyConfig;
  plex?: PlexConfig;
  jellyfin?: JellyfinConfig;
  cronjob?: string;
}

export interface EmbyOptions {
  ignoreUserId: boolean;
}

export interface TraktOptions {
  ignoreCode: boolean;
  ignoreClientSecret: boolean;
  ignoreAccessToken: boolean;
}

export interface EmbyUser {
  Id: string;
  Name: string;
}
