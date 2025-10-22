import type { Metadata } from "next";
import "./globals.css";
import { ThemeProvider } from "next-themes";
import { ReduxProvider } from "@/store/Provider";
import NextAuthSessionProvider from "./components/NextAuthSessionProvider";

export const metadata: Metadata = {
  title: "Job Gen Application",
  description: "Job Gen Application",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="font-sans antialiased">
        <NextAuthSessionProvider>
          <ReduxProvider>
            <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
              {children}
            </ThemeProvider>
          </ReduxProvider>
        </NextAuthSessionProvider>
      </body>
    </html>
  );
}
