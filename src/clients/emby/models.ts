export interface EmbyWebhook {
  Title: string;
  Date: Date;
  Event: string;
  User: {
    Name: string;
    Id: string;
  };
  Item: {
    Name: string;
    Id: string;
    ProviderIds: Record<string, string>;
    IndexNumber?: number;
    ParentIndexNumber?: number;
    Type: string;
  };
  Server: Record<string, string>;
}

export interface EmbyPlayedItemsResponse {
  PlaybackPositionTicks: number;
  PlayCount: number;
  IsFavorite: boolean;
  LastPlayedDate: string;
  Played: boolean;
}

export interface EmbyItemResponse {
  Name: string;
  Id: string;
  Type: string;
  ServerId: string;
  UserData: { Played: boolean };
  ProviderIds: Record<string, string>;
  ParentIndexNumber?: number;
  IndexNumber?: number;
  Episodes?: EmbyItemResponse[];
}

export interface EmbyItems {
  movies: EmbyItemResponse[];
  series: EmbyItemResponse[];
}

export interface EmbyUserResponse {
  Id: string;
  Name: string;
  ServerId: string;
  Prefix: string;
}
