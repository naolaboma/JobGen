import type { NextAuthOptions, User } from "next-auth";
import { JWT } from "next-auth/jwt";
import CredentialsProvider from "next-auth/providers/credentials";

// Helper function to parse JWT
const parseJwt = (token: string) => {
    try {
        return JSON.parse(Buffer.from(token.split('.')[1], 'base64').toString());
    } catch (e) {
        return null;
    }
};

// This function handles refreshing the access token
async function refreshAccessToken(token: JWT): Promise<JWT> {
    try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/refresh`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ refresh_token: token.refreshToken }),
        });

        const refreshedData = await res.json();

        if (!res.ok) {
            throw refreshedData;
        }

        const newAccessToken = refreshedData.data.access_token;
        const newRefreshToken = refreshedData.data.refresh_token || token.refreshToken; 
        const newAccessTokenExpires = parseJwt(newAccessToken).exp * 1000;

        return {
            ...token,
            accessToken: newAccessToken,
            accessTokenExpires: newAccessTokenExpires,
            refreshToken: newRefreshToken,
        };
    } catch (error) {
        console.error("RefreshAccessTokenError", error);
        return {
            ...token,
            error: "RefreshAccessTokenError",
        };
    }
}

export const authOptions: NextAuthOptions = {
    providers: [
        CredentialsProvider({
            name: "Credentials",
            credentials: {
                email: { label: "Email", type: "email" },
                password: { label: "Password", type: "password" },
            },
            async authorize(credentials) {
                if (!credentials?.email || !credentials.password) {
                    throw new Error("Please enter an email and password.");
                }

                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/login`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        email: credentials.email,
                        password: credentials.password,
                    }),
                });

                if (!res.ok) {
                    const errorData = await res.json().catch(() => ({ message: "Invalid credentials" }));
                    throw new Error(errorData.message || "Invalid credentials");
                }

                const backendResponse = await res.json();
                
                const user = backendResponse.data; 
                const accessToken = backendResponse.data.access_token;
                const refreshToken = backendResponse.data.refresh_token;

                if (user && accessToken && refreshToken) {
                    const accessTokenExpires = parseJwt(accessToken).exp * 1000;
                    return {
                        id: user.id,
                        name: user.full_name || user.username || '',
                        email: user.email,
                        accessToken,
                        refreshToken,
                        accessTokenExpires,
                    } as User & { accessToken: string; refreshToken: string; accessTokenExpires: number };
                }

                return null;
            },
        }),
    ],
    session: {
        strategy: "jwt",
    },
    callbacks: {
        async jwt({ token, user, account }) {
            // Initial sign in
            if (account && user) {
                return {
                    ...token,
                    id: user.id,
                    accessToken: (user as any).accessToken,
                    refreshToken: (user as any).refreshToken,
                    accessTokenExpires: (user as any).accessTokenExpires,
                };
            }

            // Return previous token if the access token has not expired yet
            if (Date.now() < (token.accessTokenExpires as number)) {
                return token;
            }

            // Access token has expired, try to update it
            return refreshAccessToken(token);
        },
        async session({ session, token }) {
            if (session.user) {
                session.user.id = token.id as string;
            }
            (session as any).accessToken = token.accessToken;
            (session as any).error = token.error; // Propagate error to the client

            return session;
        },
    },
    pages: {
        signIn: '/login',
        error: '/login',
    },
    secret: process.env.NEXTAUTH_SECRET,
};
