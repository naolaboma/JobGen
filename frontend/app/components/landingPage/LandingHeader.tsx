import React from "react";
import { Menu, X } from "lucide-react";
import { Button } from "@/app/components/landingPage/ui/button";
import Link from "next/link";

export function LandingHeader({
  mobileMenuOpen,
  setMobileMenuOpen,
}: {
  mobileMenuOpen: boolean;
  setMobileMenuOpen: (open: boolean) => void;
}) {
  return (
    <header className="sticky top-0 z-50 bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center">
              <div className="w-6 h-6 bg-[#44C3BB] rounded-full flex items-center justify-center">
                <div className="flex gap-1">
                  <div className="w-1 h-1 bg-white rounded-full"></div>
                  <div className="w-1 h-1 bg-white rounded-full"></div>
                </div>
              </div>
            </div>
            <span className="text-xl font-bold text-white dark:text-gray-100">
              JobGen
            </span>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-8">
            <a
              href="#features"
              className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition-colors"
            >
              Features
            </a>
            <a
              href="#pricing"
              className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition-colors"
            >
              Pricing
            </a>
            <a
              href="#partners"
              className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition-colors"
            >
              Partners
            </a>
          </nav>

          {/* Desktop CTA */}
          <div className="hidden md:flex items-center space-x-4">
            <Link href="/login" passHref>
              <Button
                variant="outline"
                className="text-black border-gray-300 hover:bg-gray-50 dark:text-gray-100 dark:border-gray-700 dark:hover:bg-gray-800"
              >
                Sign In
              </Button>
            </Link>
            <Link href="/register" passHref>
              <Button className="bg-[#26948C] hover:bg-[#1e7a73] text-white">
                Get Started
              </Button>
            </Link>
          </div>

          {/* Mobile menu button */}
          <button
            className="md:hidden p-2"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          >
            {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>

        {/* Mobile Navigation */}
        {mobileMenuOpen && (
          <div className="md:hidden py-4 border-t border-gray-100 dark:border-gray-700 bg-white dark:bg-gray-900">
            <div className="flex flex-col space-y-4">
              <a
                href="#features"
                className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition-colors"
              >
                Features
              </a>
              <a
                href="#pricing"
                className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition-colors"
              >
                Pricing
              </a>
              <a
                href="#partners"
                className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition-colors"
              >
                Partners
              </a>
              <div className="flex flex-col space-y-2 pt-4 border-t border-gray-100 dark:border-gray-700">
                <Link href="/login" passHref>
                  <Button
                    variant="outline"
                    className="w-full dark:text-gray-100 dark:border-gray-700 dark:hover:bg-gray-800"
                  >
                    Sign In
                  </Button>
                </Link>
                <Link href="/register" passHref>
                  <Button className="w-full bg-[#44C3BB] hover:bg-[#3AB5AD] text-white">
                    Get Started
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        )}
      </div>
    </header>
  );
}
