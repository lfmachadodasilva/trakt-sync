export interface TraktTypeWatched {
  title: string;
  year: number;
  ids: {
    trakt: number;
    slug: string;
    imdb: string;
  };
}

export interface TraktWatchedResponse {
  last_watched_at: Date;
  last_updated_at: Date;
  movie: TraktTypeWatched;
  show: TraktTypeWatched;
  seasons: {
    number: number;
    episodes: {
      number: number;
      last_watched_at: Date;
    }[];
  }[];
}

export interface TraktWatched {
  movies: TraktWatchedResponse[];
  shows: TraktWatchedResponse[];
}

export interface TraktAuthResponse {
  access_token: string;
  expires_in: number;
  refresh_token: string;
  token_type: string;
  scope: string;
  created_at: number;
}

export interface TraktMarkAsWatchedRequest {
  movies?: { ids: { imdb: string }; watched_at: Date }[];
  shows?: {
    ids: { imdb: string };
    seasons?: {
      watched_at: Date;
      number: number;
      episodes?: { ids: { imdb: string; watched_at: Date } }[];
    }[];
  }[];
}

export interface TraktMarkAsWatchedResponse {
  added: {
    movies: number;
    episodes: number;
  };
  updated: {
    movies: number;
    episodes: number;
  };
}
