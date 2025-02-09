import { useState } from "react";
import Image from "next/image";

import { ProfileIcon } from "@/config/icons";

export default function Avatar({ src }: { src: string | null }) {
  const [useFallback, setUseFallback] = useState(false);

  return (
    <span className="flex items-center justify-center w-[5.5rem] h-[5.5rem] bg-default-100 rounded-full overflow-hidden">
      {useFallback || !src ? (
        <ProfileIcon className="w-10 h-10" fill="currentColor" />
      ) : (
        <Image
          alt="avatar"
          height={400}
          loader={() => src}
          src={src}
          width={400}
          onError={() => setUseFallback(true)}
        />
      )}
    </span>
  );
}
