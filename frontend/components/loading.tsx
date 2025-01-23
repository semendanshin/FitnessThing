import { Spinner } from "@nextui-org/react";

export const Loading = () => {
  return (
    <div className="p4 flex flex-col flex-grow items-center justify-center gap-4">
      <Spinner />
      <p className="text-sm text-neutral-600">Загрузка...</p>
    </div>
  );
};
