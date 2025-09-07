import React from "react";

export function TestimonialsSection() {
  const testimonials = [
    {
      quote:
        "JobGen helped me land my dream remote job in just 3 weeks. The AI feedback was incredibly accurate.",
      name: "Sarah Chen",
      role: "Software Engineer",
      avatar: "SC",
    },
    {
      quote:
        "The skill gap analysis showed me exactly what I needed to learn. Now I'm working at a top tech company.",
      name: "Michael Okafor",
      role: "Product Manager",
      avatar: "MO",
    },
    {
      quote:
        "Best investment I made in my career. The interview prep feature is amazing.",
      name: "Amara Diallo",
      role: "Data Scientist",
      avatar: "AD",
    },
  ];
  return (
    <section id="partners" className="py-20 bg-white dark:bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
            Trusted by professionals worldwide
          </h2>
        </div>
        <div className="grid md:grid-cols-3 gap-8">
          {testimonials.map((testimonial, index) => (
            <div
              key={index}
              className="bg-gray-50 dark:bg-gray-800 rounded-xl p-6"
            >
              <p className="text-gray-600 dark:text-gray-300 mb-6 italic">
                "{testimonial.quote}"
              </p>
              <div className="flex items-center">
                <div className="w-12 h-12 bg-[#44C3BB] rounded-full flex items-center justify-center text-white font-semibold mr-4">
                  {testimonial.avatar}
                </div>
                <div>
                  <div className="font-semibold text-gray-900 dark:text-gray-100">
                    {testimonial.name}
                  </div>
                  <div className="text-gray-600 dark:text-gray-300 text-sm">
                    {testimonial.role}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
