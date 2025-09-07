import { CheckIcon } from "lucide-react";
import React, { JSX } from "react";
import { Card, CardContent } from "./ui/card";
import { Button } from "./ui/button";

const pricingPlans = [
  {
    title: "Free Tier",
    price: "$0",
    description:
      "Get a taste of the AI advantage and start optimizing your job search.",
    features: [
      "One comprehensive CV analysis report",
      "Basic keyword and formatting suggestions",
      "200 MB max. note size",
      "Access to 3 personalized job matches per week",
      "Standard match percentage scores",
      "Browser-based chat interface",
    ],
    buttonText: "Get Started",
    buttonVariant: "outline",
    cardStyle: "bg-white border-[#ffe492]",
    textColor: "text-[#212529]",
    priceColor: "text-[#212529]",
    buttonStyle: "bg-white border-[#ffe492] text-[#212529]",
  },
  {
    title: "Job Seeker",
    price: "$4.99",
    description:
      "Everything you need to conduct a serious, effective job search and land your next role.",
    features: [
      "Unlimited AI CV analysis and feedback",
      "Interactive CV rewrite assistance (chat-based editing)",
      "Unlimited personalized job matches",
      'Advanced match insights (e.g., "Skill Gap Analysis")',
      "Priority listing access from our web scrapers",
      "Email notifications for new, high-match jobs",
    ],
    buttonText: "Sign Up Now",
    buttonVariant: "default",
    cardStyle: "bg-[#38b2ac] shadow-[0px_4px_50px_#00000014]",
    textColor: "text-white",
    priceColor: "text-[#ffe492]",
    buttonStyle: "bg-[#7bbfb3] text-white",
    featured: true,
  },
  {
    title: "Power User",
    price: "$14.99 / month",
    description:
      "For professionals who are serious about landing top-tier global roles and fast-tracking their careers.",
    features: [
      "AI Cover Letter Generator: Draft tailored cover letters in seconds",
      "Practice with AI-generated questions based on specific job descriptions",
      "Salary Negotiation Guidance: Get data-backed advice on your offer.",
      "Exportable CV versions: Save multiple tailored versions of your optimized CV.",
      "Dedicated support: Get help faster.",
    ],
    buttonText: "Try Free for 7 Days",
    buttonVariant: "outline",
    cardStyle: "bg-white border-[#ffe492]",
    textColor: "text-[#212529]",
    priceColor: "text-[#212529]",
    buttonStyle: "bg-white border-[#ffe492] text-[#212529]",
  },
];

export const PricingSection = (): JSX.Element => {
  return (
    <section className="flex flex-col w-full items-center justify-center gap-8 md:gap-12 lg:gap-[60px] px-4 md:px-8 lg:px-[220px] py-16 md:py-24 lg:py-[140px] relative">
      <div className="inline-flex flex-col items-center justify-center gap-8 md:gap-12 lg:gap-[60px] relative flex-[0_0_auto]">
        <div className="flex flex-col w-full max-w-[1481px] items-center gap-4 md:gap-6 relative flex-[0_0_auto]">
          <img
            className="absolute w-[200px] md:w-[280px] lg:w-[319px] h-[19px] md:h-[25px] lg:h-[30px] top-[40px] md:top-[50px] lg:top-[60px] left-1/2 transform -translate-x-1/2"
            alt="Element"
            src="/element-4.png"
          />

          <h1 className="relative self-stretch mt-[-1.00px] font-header-h1 font-[number:var(--header-h1-font-weight)] text-[#212529] text-2xl md:text-4xl lg:text-[length:var(--header-h1-font-size)] text-center tracking-[var(--header-h1-letter-spacing)] leading-tight md:leading-[var(--header-h1-line-height)] [font-style:var(--header-h1-font-style)]">
            Choose Your Plan
          </h1>

          <p className="relative w-full max-w-[979px] [font-family:'Inter',Helvetica] font-normal text-[#212529] text-base md:text-lg text-center leading-6 md:leading-[18px]">
            <span className="font-bold tracking-[-0.06px] leading-6 md:leading-[30px]">
              Whether you&#39;re just starting your job search or accelerating
              your career, JobGen has a plan to fit your goals and budget.
            </span>
          </p>
        </div>

        <div className="flex w-full max-w-[1481px] items-stretch justify-center gap-4 md:gap-6 lg:gap-8 relative flex-[0_0_auto] flex-col md:flex-row">
          {pricingPlans.map((plan, index) => (
            <Card
              key={index}
              className={`flex flex-col items-start justify-center gap-4 md:gap-6 lg:gap-[25px] px-6 md:px-8 lg:px-11 ${
                plan.featured ? "py-12 md:py-16 lg:py-20" : "py-8 md:py-10"
              } relative flex-1 w-full md:min-w-[280px] lg:min-w-[300px] ${
                plan.cardStyle
              } rounded-[10px]`}
            >
              <CardContent className="flex flex-col gap-4 md:gap-6 lg:gap-[25px] p-0 w-full">
                <div className="flex flex-col items-start gap-4 md:gap-6 lg:gap-[25px] relative self-stretch w-full flex-[0_0_auto]">
                  <h3
                    className={`relative self-stretch mt-[-1.00px] font-paragraph-p1-semibold font-[number:var(--paragraph-p1-semibold-font-weight)] ${plan.textColor} text-lg md:text-xl lg:text-[length:var(--paragraph-p1-semibold-font-size)] tracking-[var(--paragraph-p1-semibold-letter-spacing)] leading-tight md:leading-[var(--paragraph-p1-semibold-line-height)] [font-style:var(--paragraph-p1-semibold-font-style)]`}
                  >
                    {plan.title}
                  </h3>

                  <div
                    className={`relative self-stretch font-header-h4 font-[number:var(--header-h4-font-weight)] ${plan.priceColor} text-2xl md:text-3xl lg:text-[length:var(--header-h4-font-size)] tracking-[var(--header-h4-letter-spacing)] leading-tight md:leading-[var(--header-h4-line-height)] [font-style:var(--header-h4-font-style)]`}
                  >
                    {plan.price}
                  </div>

                  <p
                    className={`relative self-stretch [font-family:'Inter',Helvetica] font-normal ${plan.textColor} text-sm md:text-base lg:text-lg leading-5 md:leading-6 lg:leading-[18px]`}
                  >
                    <span className="font-bold tracking-[-0.06px] leading-5 md:leading-6 lg:leading-[23px]">
                      {plan.description}
                    </span>
                  </p>
                </div>

                <div className="flex flex-col items-start gap-4 md:gap-5 lg:gap-7 relative self-stretch w-full flex-[0_0_auto]">
                  {plan.features.map((feature, featureIndex) => (
                    <div
                      key={featureIndex}
                      className="flex items-start gap-3 md:gap-4 lg:gap-[19px] relative self-stretch w-full flex-[0_0_auto]"
                    >
                      <CheckIcon
                        className={`relative w-4 h-4 md:w-5 md:h-5 mt-0.5 flex-shrink-0 ${
                          plan.featured ? "text-white" : "text-[#212529]"
                        }`}
                      />
                      <div
                        className={`relative flex-1 [font-family:'Inter',Helvetica] font-normal ${
                          plan.textColor
                        } ${
                          plan.featured
                            ? "text-sm md:text-base lg:text-lg leading-5 md:leading-6 lg:leading-[18px]"
                            : "text-sm md:text-base leading-5 md:leading-4"
                        }`}
                      >
                        <span
                          className={`font-bold ${
                            plan.featured
                              ? "tracking-[-0.06px] leading-5 md:leading-6 lg:leading-[23px]"
                              : "tracking-[-0.05px] leading-5"
                          }`}
                        >
                          {feature}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>

                <Button
                  className={`px-6 py-3 md:px-8 md:py-4 lg:px-10 ${plan.buttonStyle} rounded-lg h-auto inline-flex items-center justify-center relative flex-[0_0_auto] w-full`}
                >
                  <span className="relative w-fit mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-sm md:text-base leading-4 font-bold tracking-[-0.05px] text-center">
                    {plan.buttonText}
                  </span>
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
};
