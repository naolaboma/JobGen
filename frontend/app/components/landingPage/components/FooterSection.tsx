import { ChevronDownIcon, GlobeIcon } from "lucide-react";
import React, { JSX } from "react";
import { Button } from "./ui/button";
import { Separator } from "@radix-ui/react-separator";

export const FooterSection = (): JSX.Element => {
  const footerColumns = [
    {
      title: "Product",
      links: [
        { text: "How It Works", color: "text-[#ffe492]" },
        { text: "Features", color: "text-white" },
        { text: "Pricing", color: "text-white" },
      ],
    },
    {
      title: "Resources",
      links: [
        { text: "Blog", color: "text-white" },
        { text: "Webinars & Events", color: "text-white" },
        { text: "Success Stories", color: "text-white" },
      ],
    },
    {
      title: "Support",
      links: [
        { text: "FAQ", color: "text-white" },
        { text: "Contact Us", color: "text-white" },
        { text: "Privacy Policy", color: "text-white" },
      ],
    },
  ];

  const bottomLinks = [
    "Terms & privacy",
    "Security",
    "Status",
    "©2025 JobGen.",
  ];

  return (
    <footer className="flex flex-col w-full items-center justify-center gap-12 md:gap-20 lg:gap-[200px] pt-16 md:pt-24 lg:pt-[140px] pb-8 px-4 md:px-8 lg:px-[220px] bg-[#7bbfb3] relative">
      <div className="inline-flex flex-col justify-center gap-12 md:gap-20 lg:gap-[100px] items-center relative flex-[0_0_auto] w-full">
        <div className="w-full max-w-[1480px] gap-8 md:gap-12 lg:gap-[100px] flex-[0_0_auto] flex flex-col lg:flex-row items-start relative">
          <div className="flex-col gap-4 md:gap-6 lg:gap-[15px] flex-1 grow flex items-start relative w-full lg:w-auto">
            <div className="relative w-full lg:w-60 mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-[#f7f7ee] text-base md:text-lg leading-6 md:leading-[18px] text-center lg:text-left">
              <span className="font-bold tracking-[-0.06px] leading-[30px]">
                Connecting African tech talent to global opportunities.
              </span>
            </div>
          </div>

          {footerColumns.map((column, index) => (
            <div
              key={index}
              className="flex-col gap-3 md:gap-4 lg:gap-[15px] flex-1 grow flex items-center lg:items-start relative w-full lg:w-auto"
            >
              <div className="relative w-fit mt-[-1.00px] font-paragraph-p2-bold font-[number:var(--paragraph-p2-bold-font-weight)] text-white text-base md:text-lg lg:text-[length:var(--paragraph-p2-bold-font-size)] tracking-[var(--paragraph-p2-bold-letter-spacing)] leading-tight md:leading-[var(--paragraph-p2-bold-line-height)] [font-style:var(--paragraph-p2-bold-font-style)]">
                {column.title}
              </div>

              {column.links.map((link, linkIndex) => (
                <div
                  key={linkIndex}
                  className={`relative w-fit [font-family:'Inter',Helvetica] font-normal ${link.color} text-sm md:text-base leading-4 md:leading-5`}
                >
                  <span className="font-bold tracking-[-0.05px] leading-4 md:leading-5">
                    {link.text}
                  </span>
                </div>
              ))}
            </div>
          ))}

          <div className="inline-flex flex-col items-center lg:items-start gap-4 md:gap-6 lg:gap-[23px] relative flex-[0_0_auto] w-full lg:w-auto">
            <div className="relative w-fit mt-[-1.00px] font-header-h5 font-[number:var(--header-h5-font-weight)] text-white text-lg md:text-xl lg:text-[length:var(--header-h5-font-size)] tracking-[var(--header-h5-letter-spacing)] leading-tight md:leading-[var(--header-h5-line-height)] text-center lg:text-left [font-style:var(--header-h5-font-style)]">
              Stay Ahead in Your Job Search
            </div>

            <div className="relative w-full max-w-[259px] [font-family:'Inter',Helvetica] font-normal text-white text-sm md:text-base leading-5 md:leading-4 text-center lg:text-left">
              <span className="font-bold tracking-[-0.05px] leading-5">
                Get curated remote job leads, career tips delivered to your
                inbox.
              </span>
            </div>

            <Button className="px-6 py-4 md:px-10 md:py-5 inline-flex items-center justify-center gap-2.5 relative flex-[0_0_auto] bg-[#38b2ac] rounded-lg h-auto hover:bg-[#319795]">
              <div className="relative w-fit mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-white text-sm md:text-base leading-4">
                <span className="font-bold tracking-[-0.05px] leading-5">
                  Start today
                </span>
              </div>

              <img
                className="relative flex-[0_0_auto] mr-[-0.50px]"
                alt="Icon"
                src="/icon-17.svg"
              />
            </Button>
          </div>
        </div>

        <div className="flex w-full max-w-[1480px] items-center justify-between relative flex-[0_0_auto] flex-col lg:flex-row gap-6 lg:gap-0">
          <div className="inline-flex items-center gap-4 md:gap-8 lg:gap-[60px] relative flex-[0_0_auto] flex-wrap justify-center lg:justify-start">
            <div className="inline-flex items-center justify-center gap-1.5 relative flex-[0_0_auto]">
              <GlobeIcon className="relative flex-[0_0_auto] w-4 h-4 text-white" />

              <div className="relative w-fit mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-white text-sm md:text-base leading-4">
                <span className="font-bold tracking-[-0.05px] leading-5">
                  English
                </span>
              </div>

              <ChevronDownIcon className="relative flex-[0_0_auto] w-4 h-4 text-white" />
            </div>

            {bottomLinks.map((link, index) => (
              <div
                key={index}
                className="relative w-fit mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-white text-sm md:text-base leading-4"
              >
                <span className="font-bold tracking-[-0.05px] leading-5">
                  {link}
                </span>
              </div>
            ))}
          </div>

          <img
            className="relative flex-[0_0_auto] w-8 h-8 md:w-10 md:h-10"
            alt="Social icon"
            src="/social-icon.svg"
          />
        </div>

        <Separator className="absolute w-full max-w-[1480px] h-px top-[180px] md:top-[220px] lg:top-[248px] left-1/2 transform -translate-x-1/2 bg-white/20" />
      </div>
    </footer>
  );
};
