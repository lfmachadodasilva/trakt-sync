export interface EmbyUser {
    played: boolean;
    lastPlayedDate: string;
}

export interface EmbyProviderIds {
    imdb: string;
}

export interface EmbyItemResponse {
    id: string;
    name: string;
    serverId: string;
    type: string;
    userData: EmbyUser;
    providerIds: EmbyProviderIds;
    indexNumber: number;
    parentIndexNumber: number;
    episodes?: EmbyItemResponse[];
    parentId?: string;
}

export interface EmbyConfig {
    userId: string;
    apiKey: string;
    serverUrl: string;
}