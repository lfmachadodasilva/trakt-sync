export interface TraktSyncData {
  clientId: string;
  clientSecret: string;
  accessToken: string;
  refreshToken: string;
  code: string;
  redirectUrl: string;
}

export interface EmbySyncData {
  userId: string;
  apiKey: string;
  baseUrl: string;
}

export interface PlexSyncData {
  // TODO - Implement Plex sync data 🚧
  userId: string;
}

export interface JellyfinSyncData {
  // TODO - Implement Jellyfin sync data 🚧
  userId: string;
}

export interface SyncData {
  trakt?: TraktSyncData;
  emby?: EmbySyncData;
  plex?: PlexSyncData;
  jellyfin?: JellyfinSyncData;
}
