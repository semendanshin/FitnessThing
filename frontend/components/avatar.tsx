import { useState } from "react";
import Image from "next/image";

import { ProfileIcon } from "@/config/icons";

export default function Avatar({ src }: { src: string | null }) {
  const [useFallback, setUseFallback] = useState(false);

  return (
    <div className="relative w-[5.5rem] h-[5.5rem]">
      {useFallback || !src ? (
        <span className="flex items-center justify-center w-full h-full bg-default-100 rounded-full overflow-hidden">
          <ProfileIcon className="w-10 h-10" fill="currentColor" />
        </span>
      ) : (
        <Image
          fill
          alt="avatar"
          className="rounded-full"
          loader={() => src}
          src={src}
          style={{ objectFit: "cover" }}
          onError={() => setUseFallback(true)}
        />
      )}
    </div>
  );
}
