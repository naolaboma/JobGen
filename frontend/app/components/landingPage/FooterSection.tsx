import React from "react";

export function FooterSection() {
  return (
    <footer className="bg-gray-900 text-white py-16 dark:bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid md:grid-cols-4 gap-8">
          <div className="md:col-span-2">
            <div className="flex items-center space-x-3 mb-4">
              <div className="w-10 h-10 bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] rounded-xl flex items-center justify-center">
                <div className="w-6 h-6 bg-white rounded-lg"></div>
              </div>
              <span className="text-xl font-bold text-white dark:text-gray-100">
                JobGen
              </span>
            </div>
            <p className="text-gray-400 dark:text-gray-300 mb-6 max-w-md">
              Empowering African tech talent to access global opportunities
              through AI-powered career tools.
            </p>
            <div className="flex space-x-4">
              <div className="w-10 h-10 bg-gray-800 dark:bg-gray-700 rounded-lg flex items-center justify-center hover:bg-gray-700 dark:hover:bg-gray-600 cursor-pointer">
                <span className="text-sm text-white">G</span>
              </div>
              <div className="w-10 h-10 bg-gray-800 dark:bg-gray-700 rounded-lg flex items-center justify-center hover:bg-gray-700 dark:hover:bg-gray-600 cursor-pointer">
                <span className="text-sm text-white">in</span>
              </div>
              <div className="w-10 h-10 bg-gray-800 dark:bg-gray-700 rounded-lg flex items-center justify-center hover:bg-gray-700 dark:hover:bg-gray-600 cursor-pointer">
                <span className="text-sm text-white">f</span>
              </div>
            </div>
          </div>
          <div>
            <h3 className="font-semibold mb-4 text-white dark:text-gray-100">
              Product
            </h3>
            <ul className="space-y-2 text-gray-400 dark:text-gray-300">
              <li>
                <a
                  href="#"
                  className="hover:text-white dark:hover:text-[#44C3BB] transition-colors"
                >
                  Features
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Pricing
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  API
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Integrations
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="font-semibold mb-4 text-white dark:text-gray-100">
              Support
            </h3>
            <ul className="space-y-2 text-gray-400 dark:text-gray-300">
              <li>
                <a
                  href="#"
                  className="hover:text-white dark:hover:text-[#44C3BB] transition-colors"
                >
                  Help Center
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Contact Us
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Privacy Policy
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Terms of Service
                </a>
              </li>
            </ul>
          </div>
        </div>
        <div className="border-t border-gray-800 dark:border-gray-700 mt-12 pt-8 text-center text-gray-400 dark:text-gray-300">
          <p>&copy; 2025 JobGen. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
}
