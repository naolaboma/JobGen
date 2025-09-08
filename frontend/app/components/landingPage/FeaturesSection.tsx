import React from "react";

export function FeaturesSection() {
  const features = [
    {
      title: "AI CV Analysis",
      description:
        "Get instant feedback on your resume with AI-powered analysis and optimization suggestions.",
      icon: "📄",
    },
    {
      title: "Smart Job Matching",
      description:
        "Discover opportunities that perfectly match your skills and career goals.",
      icon: "🎯",
    },
    {
      title: "Interview Prep",
      description:
        "Practice with AI-generated questions tailored to specific job descriptions.",
      icon: "💬",
    },
    {
      title: "Skill Gap Analysis",
      description:
        "Understand exactly what skills you need to develop for your target roles.",
      icon: "📊",
    },
    {
      title: "Cover Letter Generator",
      description:
        "Create compelling cover letters in seconds with our AI writing assistant.",
      icon: "✍️",
    },
    {
      title: "Salary Insights",
      description:
        "Get data-backed advice for salary negotiations and career planning.",
      icon: "💰",
    },
  ];
  return (
    <section id="features" className="py-20 bg-white dark:bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
            Everything you need to succeed
          </h2>
          <p className="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
            Our AI-powered platform provides comprehensive tools to optimize
            your job search and land your dream role.
          </p>
        </div>
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, index) => (
            <div
              key={index}
              className="bg-gray-50 dark:bg-gray-800 rounded-xl p-6 hover:shadow-lg transition-shadow"
            >
              <div className="text-4xl mb-4">{feature.icon}</div>
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-3">
                {feature.title}
              </h3>
              <p className="text-gray-600 dark:text-gray-300">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
