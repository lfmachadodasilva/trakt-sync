export interface TraktConfig {
    clientId: string;
    clientSecret: string;
    accessToken: string;
    refreshToken: string;
    code: string;
    redirectUrl: string;
    expiresIn: number;
    createdAt: number;
}