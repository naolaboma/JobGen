import React from "react";
import { Button } from "@/app/components/landingPage/ui/button";
import Link from "next/link";

export function HeroSection() {
  return (
    <section className="relative overflow-hidden bg-gradient-to-br from-teal-50 via-white to-blue-50 dark:from-gray-900 dark:via-gray-900 dark:to-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 lg:py-32">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          <div className="text-center lg:text-left">
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-gray-900 dark:text-gray-100 leading-tight mb-6">
              Land Your Dream
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-teal-600 to-blue-600">
                {" "}
                Remote Tech Job
              </span>
            </h1>
            <p className="text-xl text-gray-600 dark:text-gray-300 mb-8 leading-relaxed">
              AI-powered career coaching for African tech talent. Get
              personalized CV feedback and discover remote opportunities that
              match your skills.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
              <Link href="/register" passHref>
                <Button
                  size="lg"
                  className="bg-[#44C3BB] hover:bg-[#3AB5AD] text-white px-8 py-4 text-lg"
                >
                  Start Free Trial
                </Button>
              </Link>
              <Button
                size="lg"
                variant="outline"
                className="px-8 py-4 text-lg dark:border-gray-100 dark:text-[#44C3BB] dark:hover:bg-gray-100 dark:hover:text-[#44C3BB]"
              >
                Watch Demo
              </Button>
            </div>
          </div>
          <div className="relative">
            <div className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-800 dark:to-gray-700 rounded-2xl p-8 shadow-2xl">
              <div className="bg-white dark:bg-gray-900 rounded-xl p-6 space-y-4">
                <div className="flex items-center space-x-3">
                  <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                  <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                  <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                </div>
                <div className="space-y-3">
                  <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-3/4"></div>
                  <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-1/2"></div>
                  <div className="h-4 bg-teal-200 dark:bg-[#44C3BB]/30 rounded w-2/3"></div>
                  <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-full"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
