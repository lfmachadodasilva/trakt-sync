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
  Episodes: EmbyItemResponse[];
}

export interface EmbyItems {
  movies: EmbyItemResponse[];
  series: EmbyItemResponse[];
}
