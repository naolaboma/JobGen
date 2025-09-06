import type { NextAuthOptions } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";

export const authOptions: NextAuthOptions = {
    providers: [
        CredentialsProvider({
            name: "Credentials", // Consistent naming
            credentials: {
                email: { label: "Email", type: "email" },
                password: { label: "Password", type: "password" },
            },
            async authorize(credentials) {
                if (!credentials?.email || !credentials.password) {
                    throw new Error("Please enter an email and password."); // More specific error
                }

                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/login`, { // Use NEXT_PUBLIC_API_URL
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        email: credentials.email,
                        password: credentials.password,
                    }),
                });

                if (!res.ok) {
                    const errorData = await res.json().catch(() => ({ message: "Invalid credentials" }));
                    throw new Error(errorData.message || "Invalid credentials"); // Robust error handling
                }

                const backendResponse = await res.json();
                const user = backendResponse.data; // Assuming 'data' contains the user object

                if (user) {
                    // Harmonize user object for NextAuth
                    return {
                        id: user.id,
                        name: user.full_name || user.username || '', // Use full_name or username for name, default to empty string
                        email: user.email,
                        accessToken: backendResponse.access_token, // Assuming tokens are directly in backendResponse
                        refreshToken: backendResponse.refresh_token,
                    };
                }

                return null;
            },
        }),
    ],
    session: {
        strategy: "jwt",
    },
    callbacks: {
        async jwt({ token, user }) {
            if (user) {
                token.id = user.id;
                token.accessToken = user.accessToken;
                token.refreshToken = user.refreshToken;
            }
            return token;
        },
        async session({ session, token }) {
            if (session.user) {
                session.user.id = token.id as string;
            }
            (session as any).accessToken = token.accessToken; // Ensure accessToken is on session
            (session as any).refreshToken = token.refreshToken; // Ensure refreshToken is on session
            return session;
        },
    },
    pages: {
        signIn: '/login', // Consistent with file structure
        error: '/login', // Add error page for consistency
    },
    secret: process.env.NEXTAUTH_SECRET, // Add secret
};
