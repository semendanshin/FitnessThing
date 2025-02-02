import { Button } from "@nextui-org/button";
import { Input } from "@nextui-org/input";
import clsx from "clsx";

export function IncrementButtons({
  value,
  setValue,
  isSubtract,
  radius,
}: {
  value: number;
  setValue: (value: number) => void;
  isSubtract?: boolean;
  radius?: "sm" | "md" | "lg";
}) {
  return (
    <div className="flex flex-col h-full">
      <Button
        isIconOnly
        className="min-w-fit w-fit p-3 h-full"
        radius={radius}
        onPress={() => {
          if (value > 0 && isSubtract) {
            setValue(value - 1);

            return;
          }
          if (!isSubtract) {
            setValue(value + 1);
          }
        }}
      >
        {isSubtract ? "-" : "+"}
      </Button>
    </div>
  );
}

export function InputWithIncrement({
  value,
  setValue,
  label,
  placeholder,
  type,
  className,
  size,
}: {
  value: number;
  setValue: (value: number) => void;
  label: string;
  placeholder: string;
  type: string;
  className?: string;
  size?: "sm" | "md" | "lg";
}) {
  return (
    <>
      <p>{label}</p>
      <div className="flex flex-row gap-2 items-center h-full">
        <IncrementButtons
          isSubtract
          radius={size}
          setValue={setValue}
          value={value}
        />
        <Input
          isRequired
          className={clsx("p-0 w-full h-full", className)}
          classNames={{ inputWrapper: "h-full max-h-full" }}
          placeholder={placeholder}
          size={size}
          type={type}
          value={String(value)}
          onValueChange={(value) => setValue(Number(value))}
        />
        <IncrementButtons radius={size} setValue={setValue} value={value} />
      </div>
    </>
  );
}
