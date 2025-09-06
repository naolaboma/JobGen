import React, { JSX } from "react";
import { Card, CardContent } from "./ui/card";
import { Avatar, AvatarImage } from "@radix-ui/react-avatar";

export const TestimonialsSection = (): JSX.Element => {
  const testimonials = [
    {
      id: 1,
      quote:
        "Within two weeks of using JobGen's suggestions and job matches, I got interview calls from two European startups.",
      name: "Kalkidan T.",
      title: "Software Developer, Ethiopia",
      avatar: "/avater.png",
      bgColor: "bg-white",
      textColor: "text-[#212529]",
      borderColor: "border-[#212529]",
      nameColor: "text-[#212529]",
      titleColor: "text-[#212529]",
    },
    {
      id: 2,
      quote:
        "Within two weeks of using JobGen's suggestions and job matches, I got interview calls from two European startups.",
      name: "Robel W.",
      title: "Software Developer, Ethiopia",
      avatar: "/avater-1.png",
      bgColor: "bg-[#38b2ac]",
      textColor: "text-white",
      borderColor: "border-white",
      nameColor: "text-[#212529]",
      titleColor: "text-white",
    },
    {
      id: 3,
      quote:
        "Within two weeks of using JobGen's suggestions and job matches, I got interview calls from two European startups.",
      name: "Yabsira Z.",
      title: "Software Developer, Ethiopia",
      avatar: "/avater-2.png",
      bgColor: "bg-[#38b2ac]",
      textColor: "text-white",
      borderColor: "border-white",
      nameColor: "text-[#043873]",
      titleColor: "text-white",
    },
  ];

  return (
    <section className="flex flex-col w-full items-center justify-center gap-8 md:gap-12 lg:gap-[60px] px-4 md:px-8 lg:px-[220px] py-16 md:py-24 lg:py-[140px] relative">
      <img
        className="absolute w-[200px] md:w-[230px] lg:w-[258px] h-8 md:h-10 lg:h-12 top-[120px] md:top-[150px] lg:top-[172px] right-4 md:right-8 lg:right-[220px] hidden lg:block"
        alt="Group"
        src="/group.png"
      />

      <div className="inline-flex flex-col items-center justify-center gap-8 md:gap-12 lg:gap-[60px] relative flex-[0_0_auto]">
        <h2 className="relative w-full max-w-[1479px] mt-[-1.00px] [font-family:'Inter',Helvetica] font-bold text-[#212529] text-2xl md:text-4xl lg:text-[70px] text-center tracking-[0] leading-tight md:leading-[84px] px-4">
          Trusted by Tech Professionals Across Africa
        </h2>

        <div className="flex w-full max-w-[1479px] gap-4 md:gap-6 lg:gap-8 items-stretch relative flex-[0_0_auto] flex-col lg:flex-row">
          {testimonials.map((testimonial) => (
            <Card
              key={testimonial.id}
              className={`flex flex-col items-start gap-8 md:gap-12 lg:gap-[60px] px-6 md:px-8 lg:px-10 py-8 md:py-12 lg:py-[60px] relative flex-1 grow ${
                testimonial.bgColor
              } rounded-[10px] ${
                testimonial.id === 1 ? "shadow-[15px_10px_50px_#0000001a]" : ""
              } border-0`}
            >
              <CardContent className="p-0 w-full">
                <div
                  className={`flex flex-col items-start gap-6 md:gap-8 pt-0 pb-6 md:pb-8 lg:pb-10 px-0 relative self-stretch w-full flex-[0_0_auto] border-b border-solid ${testimonial.borderColor}`}
                >
                  <img
                    className="relative flex-[0_0_auto] w-8 h-8 md:w-10 md:h-10"
                    alt="Quote"
                    src="/quote.svg"
                  />

                  <div
                    className={`relative self-stretch [font-family:'Inter',Helvetica] font-normal ${testimonial.textColor} text-sm md:text-base lg:text-lg leading-5 md:leading-6 lg:leading-[18px]`}
                  >
                    <span className="font-bold tracking-[-0.06px] leading-5 md:leading-6 lg:leading-[30px]">
                      {testimonial.quote}
                    </span>
                  </div>
                </div>

                <div className="flex w-full items-center gap-4 md:gap-6 lg:gap-[42px] relative flex-[0_0_auto] mt-6 md:mt-8 lg:mt-[60px]">
                  <Avatar className="relative w-12 h-12 md:w-16 md:h-16 lg:w-[95px] lg:h-[95px] flex-shrink-0">
                    <AvatarImage
                      src={testimonial.avatar}
                      alt="Avatar"
                      className="object-cover"
                    />
                  </Avatar>

                  <div className="flex flex-col items-start gap-2 md:gap-3 lg:gap-[15px] relative flex-1">
                    <div
                      className={`relative w-full mt-[-1.00px] font-paragraph-p1-semibold font-[number:var(--paragraph-p1-semibold-font-weight)] ${testimonial.nameColor} text-base md:text-lg lg:text-[length:var(--paragraph-p1-semibold-font-size)] tracking-[var(--paragraph-p1-semibold-letter-spacing)] leading-tight md:leading-[var(--paragraph-p1-semibold-line-height)] [font-style:var(--paragraph-p1-semibold-font-style)]`}
                    >
                      {testimonial.name}
                    </div>

                    <div
                      className={`relative w-full [font-family:'Inter',Helvetica] font-normal ${testimonial.titleColor} text-sm md:text-base leading-4 md:leading-5`}
                    >
                      <span className="font-bold tracking-[-0.05px] leading-4 md:leading-5">
                        {testimonial.title}
                      </span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
};
