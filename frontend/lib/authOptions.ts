import type { NextAuthOptions, User } from "next-auth";
import { JWT } from "next-auth/jwt";
import CredentialsProvider from "next-auth/providers/credentials";
import GoogleProvider from "next-auth/providers/google";
import GitHubProvider from "next-auth/providers/github";
import { apiUrl } from "@/lib/api";

// Helper function to parse JWT
const parseJwt = (token: string) => {
  try {
    return JSON.parse(Buffer.from(token.split(".")[1], "base64").toString());
  } catch (e) {
    return null;
  }
};

// This function handles refreshing the access token
async function refreshAccessToken(token: JWT): Promise<JWT> {
  try {
    const res = await fetch(apiUrl("/api/v1/auth/refresh"), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: token.refreshToken }),
    });

    const refreshedData = await res.json();

    if (!res.ok) {
      throw refreshedData;
    }

    const newAccessToken = refreshedData.data.access_token;
    const newRefreshToken =
      refreshedData.data.refresh_token || token.refreshToken;
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
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID as string,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET as string,
    }),
    GitHubProvider({
      clientId: process.env.GITHUB_CLIENT_ID as string,
      clientSecret: process.env.GITHUB_CLIENT_SECRET as string,
    }),
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

        const res = await fetch(apiUrl("/api/v1/auth/login"), {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            email: credentials.email,
            password: credentials.password,
          }),
        });

        if (!res.ok) {
          const errorData = await res
            .json()
            .catch(() => ({ message: "Invalid credentials" }));
          throw new Error(errorData.message || "Invalid credentials");
        }

        const backendResponse = await res.json();
        const accessToken = backendResponse.data?.access_token;
        const refreshToken = backendResponse.data?.refresh_token; // now provided by backend JSON

        if (accessToken) {
          const accessTokenPayload = parseJwt(accessToken);
          const accessTokenExpires = accessTokenPayload?.exp
            ? accessTokenPayload.exp * 1000
            : Date.now() + 60 * 60 * 1000;

          const user: User = {
            id: credentials.email,
            email: credentials.email,
            name: credentials.email,
          };

          return {
            ...user,
            accessToken,
            refreshToken, // persist refresh token into JWT
            accessTokenExpires,
          } as User & {
            accessToken: string;
            refreshToken?: string;
            accessTokenExpires: number;
          };
        }

        return null;
      },
    }),
  ],
  session: {
    strategy: "jwt",
    maxAge: 24 * 60 * 60, // 1 day
    updateAge: 60 * 60, // 1 hour
  },
  callbacks: {
    async jwt({ token, user, account }) {
      if (account && user) {
        return {
          ...token,
          id: (user as any).id,
          accessToken: (user as any).accessToken,
          refreshToken: (user as any).refreshToken, // retain refresh token on initial sign in
          accessTokenExpires: (user as any).accessTokenExpires,
        };
      }

      if (Date.now() < (token.accessTokenExpires as number)) {
        return token;
      }

      return refreshAccessToken(token as any);
    },
    async session({ session, token }) {
      if (session.user) {
        (session.user as any).id = token.id as string;
      }
      (session as any).accessToken = token.accessToken;
      (session as any).error = token.error;
      return session;
    },
  },
  pages: {
    signIn: "/login",
    error: "/login",
  },
  secret: process.env.NEXTAUTH_SECRET,
};
