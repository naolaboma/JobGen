import React, { JSX } from "react";
import { Button } from "./ui/button";

const contentBlocks = [
  {
    id: 1,
    title: "Your AI-Powered Assistant for a Global Career",
    description:
      "JobGen cuts through the noise. We combine AI-driven CV optimization with intelligent job matching to connect you directly with remote tech roles that are a perfect fit for your unique profile.",
    buttonText: "Get Started",
    hasImage: false,
    imageFirst: false,
    decorativeElements: [
      {
        src: "/element-1.png",
        className: "absolute w-[679px] h-[87px] top-[149px] left-0",
      },
    ],
  },
  {
    id: 2,
    title: "Reclaim Your Time",
    description:
      "No more spending hours scrolling through irrelevant job posts. Get a shortlist of high-quality, matched opportunities delivered to you in seconds.",
    buttonText: "Try it now",
    hasImage: true,
    imageFirst: true,
    imageSrc: "/work-together-image.png",
    decorativeElements: [
      {
        src: "/background.png",
        className: "absolute w-[441px] h-[449px] top-[-706px] left-[-1030px]",
      },
      {
        src: "/element-2.png",
        className: "absolute w-[298px] h-[29px] top-14 left-[209px]",
      },
    ],
  },
];

export const WorkManagementSection = (): JSX.Element => {
  return (
    <section className="flex flex-col w-full items-center justify-center gap-12 md:gap-20 lg:gap-[100px] px-4 md:px-8 lg:px-[220px] py-16 md:py-24 lg:py-[140px] relative">
      <div className="inline-flex flex-col items-center justify-center gap-12 md:gap-20 lg:gap-[100px] relative flex-[0_0_auto]">
        {contentBlocks.map((block) => (
          <div
            key={block.id}
            className={`flex w-full max-w-[1480px] items-center justify-center gap-8 md:gap-12 lg:gap-[60px] relative flex-[0_0_auto] ${
              block.imageFirst
                ? "flex-col lg:flex-row"
                : "flex-col lg:flex-row-reverse"
            }`}
          >
            {block.imageFirst && block.hasImage && (
              <img
                className="relative w-full max-w-[500px] lg:max-w-[710px] h-auto lg:h-[661px] object-contain"
                alt="Work together image"
                src={block.imageSrc}
              />
            )}

            <div className="flex flex-col items-start justify-center gap-8 md:gap-12 lg:gap-[60px] relative flex-1 grow">
              <div className="flex flex-col items-start gap-4 md:gap-6 relative self-stretch w-full flex-[0_0_auto]">
                {block.decorativeElements.map((element, index) => (
                  <img
                    key={`decorative-${block.id}-${index}`}
                    className={`${element.className} hidden lg:block`}
                    alt="Element"
                    src={element.src}
                  />
                ))}

                <h2 className="relative self-stretch mt-[-1.00px] font-header-h1 font-[number:var(--header-h1-font-weight)] text-[#212529] text-2xl md:text-4xl lg:text-[length:var(--header-h1-font-size)] tracking-[var(--header-h1-letter-spacing)] leading-tight md:leading-[var(--header-h1-line-height)] [font-style:var(--header-h1-font-style)]">
                  {block.title}
                </h2>

                <p className="relative self-stretch [font-family:'Inter',Helvetica] font-normal text-[#212529] text-base md:text-lg leading-6 md:leading-[18px]">
                  <span className="font-bold tracking-[-0.06px] leading-6 md:leading-[30px]">
                    {block.description}
                  </span>
                </p>
              </div>

              <Button className="inline-flex items-center justify-center gap-2.5 px-6 py-4 md:px-10 md:py-5 relative flex-[0_0_auto] bg-[#7bbfb3] rounded-md h-auto hover:bg-[#6ba89c]">
                <span className="relative w-fit mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-white text-base md:text-lg leading-[18px]">
                  <span className="font-bold tracking-[-0.06px] leading-[23px]">
                    {block.buttonText}
                  </span>
                </span>

                <img
                  className="relative flex-[0_0_auto] mr-[-0.50px]"
                  alt="Icon"
                  src="/icon-17.svg"
                />
              </Button>
            </div>

            {!block.imageFirst && !block.hasImage && (
              <div className="relative w-full max-w-[500px] lg:max-w-[748px] h-[300px] md:h-[400px] lg:h-[547px] bg-primary-100 rounded-lg" />
            )}
          </div>
        ))}
      </div>
    </section>
  );
};
