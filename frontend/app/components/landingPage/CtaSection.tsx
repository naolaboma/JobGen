import React from "react";
import { Button } from "@/app/components/landingPage/ui/button";
import Link from "next/link";

export function CtaSection() {
  return (
    <section className="py-20 bg-gradient-to-r from-teal-600 to-blue-600">
      <div className="max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8">
        <h2 className="text-3xl sm:text-4xl font-bold text-white mb-6">
          Ready to accelerate your career?
        </h2>
        <p className="text-xl text-teal-100 mb-8">
          Join thousands of professionals who have transformed their careers
          with JobGen.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link href="/register" passHref>
            <Button
              size="lg"
              className="bg-white text-teal-600 hover:bg-gray-100 px-8 py-4 text-lg"
            >
              Start Free Trial
            </Button>
          </Link>
          <Button
            size="lg"
            variant="outline"
            className="border-white text-white hover:bg-white hover:text-teal-600 px-8 py-4 text-lg"
          >
            Schedule Demo
          </Button>
        </div>
      </div>
    </section>
  );
}
