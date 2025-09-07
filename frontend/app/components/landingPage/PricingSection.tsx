import React from "react";
import { Button } from "@/app/components/landingPage/ui/button";

export function PricingSection() {
  const plans = [
    {
      name: "Free",
      price: "$0",
      period: "forever",
      features: [
        "1 CV analysis per month",
        "Basic job matching",
        "Community support",
        "Basic templates",
      ],
      cta: "Get Started",
      popular: false,
    },
    {
      name: "Job Seeker",
      price: "$9",
      period: "per month",
      features: [
        "Unlimited CV analysis",
        "Advanced job matching",
        "Interview preparation",
        "Priority support",
        "Premium templates",
        "Skill gap analysis",
      ],
      cta: "Start Free Trial",
      popular: true,
    },
    {
      name: "Power User",
      price: "$19",
      period: "per month",
      features: [
        "Everything in Job Seeker",
        "AI cover letter generator",
        "Salary negotiation guide",
        "1-on-1 career coaching",
        "Custom branding",
        "API access",
      ],
      cta: "Contact Sales",
      popular: false,
    },
  ];
  return (
    <section id="pricing" className="py-20 bg-gray-50 dark:bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
            Choose your plan
          </h2>
          <p className="text-xl text-gray-600 dark:text-gray-300">
            Start free and upgrade as you grow your career
          </p>
        </div>
        <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
          {plans.map((plan, index) => (
            <div
              key={index}
              className={`bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg ${
                plan.popular
                  ? "ring-2 ring-[#44C3BB] dark:ring-[#44C3BB] relative"
                  : ""
              }`}
            >
              {plan.popular && (
                <div className="absolute -top-4 left-1/2 transform -translate-x-1/2">
                  <span className="bg-teal-500 text-white px-4 py-2 rounded-full text-sm font-medium">
                    Most Popular
                  </span>
                </div>
              )}
              <div className="text-center mb-8">
                <h3 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">
                  {plan.name}
                </h3>
                <div className="mb-4">
                  <span className="text-4xl font-bold text-gray-900 dark:text-gray-100">
                    {plan.price}
                  </span>
                  <span className="text-gray-600 dark:text-gray-300">
                    /{plan.period}
                  </span>
                </div>
              </div>
              <ul className="space-y-4 mb-8">
                {plan.features.map((feature, featureIndex) => (
                  <li key={featureIndex} className="flex items-center">
                    <svg
                      className="w-5 h-5 text-[#44C3BB] mr-3"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                        clipRule="evenodd"
                      />
                    </svg>
                    <span className="text-gray-600 dark:text-gray-300">
                      {feature}
                    </span>
                  </li>
                ))}
              </ul>
              <Button
                className={`w-full ${
                  plan.popular
                    ? "bg-[#44C3BB] hover:bg-[#3AB5AD] text-white"
                    : "bg-gray-100 hover:bg-gray-200 text-gray-900 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-gray-100"
                }`}
              >
                {plan.cta}
              </Button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
