"use client";

import { useState, useEffect } from "react";
import { useTheme } from "next-themes";
import {
  ArrowLeft,
  Settings,
  User,
  Shield,
  HelpCircle,
  SlidersHorizontal,
  Globe,
  Moon,
  Mail,
  LogOut,
  ExternalLink,
  Bell,
  FileText,
  Cookie,
  BookOpen,
  ChevronRight,
  KeyRound,
  Menu,
  X,
  ShieldCheck,
  CircleUser,
  BadgeInfo,
} from "lucide-react";
import { Button } from "@/app/components/ui/button";
import { Card, CardContent, CardHeader } from "@/app/components/customCard";
import { Switch } from "@/app/components/ui/switch";
import Link from "next/link";
import Footer from "@/app/components/footer";

export default function SettingsPage() {
  const { theme, setTheme } = useTheme();
  const isDark = theme === "dark";
  const [emailUpdates, setEmailUpdates] = useState(true);
  const [language, setLanguage] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col">
      {/* Header */}
      <header className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center">
                <div className="w-6 h-6 bg-[#44C3BB] rounded-full flex items-center justify-center">
                  <div className="flex gap-1">
                    <div className="w-1 h-1 bg-white rounded-full"></div>
                    <div className="w-1 h-1 bg-white rounded-full"></div>
                  </div>
                </div>
              </div>
              <span className="text-xl font-semibold text-white dark:text-gray-100">
                JobGen
              </span>
            </div>

            {/* Desktop Navigation */}
            <nav className="hidden md:flex items-center gap-8">
              <Link
                href="/blog"
                className="hover:text-gray-200 transition-colors"
              >
                Blog
              </Link>
              <Link
                href="/about"
                className="hover:text-gray-200 transition-colors"
              >
                About
              </Link>
              <Link
                href="/contact"
                className="hover:text-gray-200 transition-colors"
              >
                Contact
              </Link>
            </nav>

            {/* Mobile Menu Button */}
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="md:hidden p-2 rounded-md hover:bg-white/10 transition-colors cursor-pointer"
            >
              {mobileMenuOpen ? (
                <X className="w-6 h-6" />
              ) : (
                <Menu className="w-6 h-6" />
              )}
            </button>
          </div>

          {/* Mobile Navigation */}
          {mobileMenuOpen && (
            <div className="md:hidden py-4 border-t border-white/20">
              <div className="flex flex-col gap-4">
                <Link
                  href="/blog"
                  onClick={() => setMobileMenuOpen(false)}
                  className="hover:text-gray-200 transition-colors"
                >
                  Blog
                </Link>
                <Link
                  href="/about"
                  onClick={() => setMobileMenuOpen(false)}
                  className="hover:text-gray-200 transition-colors"
                >
                  About
                </Link>
                <Link
                  href="/contact"
                  onClick={() => setMobileMenuOpen(false)}
                  className="hover:text-gray-200 transition-colors"
                >
                  Contact
                </Link>
              </div>
            </div>
          )}
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto my-2 px-6 py-8 flex-1">
        {/* Back to Dashboard */}
        <Link
          href="/dashboard"
          className="inline-flex items-center gap-2 text-gray-600 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white mb-8 transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          Back to Dashboard
        </Link>

        {/* Page Header */}
        <div className="text-center mb-12">
          <div className="flex items-center justify-center gap-3 mb-4">
            <Settings className="w-8 h-8 text-[#44C3BB] dark:text-[#44C3BB]" />
            <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
              Settings & Help
            </h1>
          </div>
          <p className="text-gray-600 dark:text-gray-300">
            Manage your account
          </p>
        </div>

        {/* Settings Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-12">
          <Card className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl h-full shadow-sm -m-px">
            <CardHeader className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white dark:text-gray-100 p-0 rounded-t-xl">
              <div className="flex items-center gap-3 p-6">
                <User className="w-6 h-6" />
                <h2 className="text-xl font-semibold">Account Settings</h2>
              </div>
            </CardHeader>
            <CardContent className="p-0">
              <div className="p-6 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <Link href={"/profile"}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <CircleUser className="w-5 h-5 text-[#44C3BB]" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                          Profile Information
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          Update your personal details
                        </p>
                      </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-gray-400 dark:text-gray-300" />
                  </div>
                </Link>
              </div>
              <div className="p-6 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <Link href={"/notifications"}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <Bell className="w-5 h-5 text-[#44C3BB]" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                          Notifications
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          Manage your notification preferences
                        </p>
                      </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-gray-400 dark:text-gray-300" />
                  </div>
                </Link>
              </div>
              <div className="p-6 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <Link href={"/security"}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <KeyRound className="w-5 h-5 text-[#44C3BB]" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                          Password & Security
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          Change password and security settings
                        </p>
                      </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-gray-400 dark:text-gray-300" />
                  </div>
                </Link>
              </div>
            </CardContent>
          </Card>

          <Card className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl h-full shadow-sm -m-px">
            <CardHeader className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white dark:text-gray-100 p-0 rounded-t-xl">
              <div className="flex items-center gap-3 p-6">
                <Shield className="w-6 h-6" />
                <h2 className="text-xl font-semibold">Privacy & Legal</h2>
              </div>
            </CardHeader>
            <CardContent className="p-0">
              <div className="p-6 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <Link href={"/privacy"}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <ShieldCheck className="w-5 h-5 text-[#44C3BB]" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                          Privacy Policy
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          Read our privacy policy
                        </p>
                      </div>
                    </div>
                    <ExternalLink className="w-5 h-5 text-gray-400 dark:text-gray-300" />
                  </div>
                </Link>
              </div>
              <div className="p-6 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <Link href={"/terms"}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <FileText className="w-5 h-5 text-[#44C3BB]" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                          Terms of Service
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          View terms and conditions
                        </p>
                      </div>
                    </div>
                    <ExternalLink className="w-5 h-5 text-gray-400 dark:text-gray-300" />
                  </div>
                </Link>
              </div>
              <div className="p-6 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <Link href={"/cookies"}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <Cookie className="w-5 h-5 text-[#44C3BB]" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                          Cookie Preferences
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          Manage cookie settings
                        </p>
                      </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-gray-400 dark:text-gray-300" />
                  </div>
                </Link>
              </div>
            </CardContent>
          </Card>

          <Card className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl h-full shadow-sm -m-px">
            <CardHeader className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white dark:text-gray-100 p-0 rounded-t-xl">
              <div className="flex items-center gap-3 p-6">
                <HelpCircle className="w-6 h-6" />
                <h2 className="text-xl font-semibold">Support & Help</h2>
              </div>
            </CardHeader>
            <CardContent className="p-0">
              <div className="p-6 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <Link href={"/support"}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <HelpCircle className="w-5 h-5 text-[#44C3BB]" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                          Contact Support
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          Get help from our support team
                        </p>
                      </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-gray-400 dark:text-gray-300" />
                  </div>
                </Link>
              </div>
              <div className="p-6 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <Link href={"/faq"}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <BadgeInfo className="w-5 h-5 text-[#44C3BB]" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                          FAQ
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          Find answers to common questions
                        </p>
                      </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-gray-400 dark:text-gray-300" />
                  </div>
                </Link>
              </div>
              <div className="p-6 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <Link href={"/user-guide"}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <BookOpen className="w-5 h-5 text-[#44C3BB]" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                          User Guide
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          Learn how to use JobGen
                        </p>
                      </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-gray-400 dark:text-gray-300" />
                  </div>
                </Link>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Preferences Card */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <Card className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl h-full shadow-sm -m-px">
            <CardHeader className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white dark:text-gray-100 p-0 rounded-t-xl">
              <div className="flex items-center gap-3 p-6">
                <SlidersHorizontal className="w-6 h-6" />
                <h2 className="text-xl font-semibold">Preferences</h2>
              </div>
            </CardHeader>
            <CardContent className="p-0">
              <div className="p-6 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <Globe className="w-5 h-5 text-[#44C3BB]" />
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                        Language
                      </h3>
                      <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                        English (US)
                      </p>
                    </div>
                  </div>
                  <Switch
                    defaultChecked={language}
                    onCheckedChange={setLanguage}
                    className="data-[state=checked]:bg-[#44C3BB] data-[state=unchecked]:bg-gray-200 dark:data-[state=unchecked]:bg-gray-700 cursor-pointer"
                  />
                </div>
              </div>
              <div className="p-6 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <Moon className="w-5 h-5 text-[#44C3BB]" />
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                        Dark Mode
                      </h3>
                      <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                        Switch to dark theme
                      </p>
                    </div>
                  </div>
                  {mounted && (
                    <Switch
                      defaultChecked={isDark}
                      onCheckedChange={(checked: boolean) =>
                        setTheme(checked ? "dark" : "light")
                      }
                      className="data-[state=checked]:bg-[#44C3BB] data-[state=unchecked]:bg-gray-200 dark:data-[state=unchecked]:bg-gray-700 cursor-pointer"
                    />
                  )}
                </div>
              </div>
              <div className="p-6 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors cursor-pointer">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <Mail className="w-5 h-5 text-[#44C3BB]" />
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-base">
                        Email Updates
                      </h3>
                      <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                        Receive job alerts via email
                      </p>
                    </div>
                  </div>
                  <Switch
                    defaultChecked={emailUpdates}
                    onCheckedChange={setEmailUpdates}
                    className="data-[state=checked]:bg-[#44C3BB] data-[state=unchecked]:bg-gray-200 dark:data-[state=unchecked]:bg-gray-700 cursor-pointer"
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Logout Button */}
        <div className="mt-12 flex justify-center">
          <Button
            variant="destructive"
            className="bg-red-100 dark:bg-red-900 text-red-600 dark:text-red-300 hover:bg-red-200 dark:hover:bg-red-800 border-0 cursor-pointer"
          >
            <LogOut className="w-4 h-4 mr-2" />
            Logout
          </Button>
        </div>
      </main>

      <Footer />
    </div>
  );
}
