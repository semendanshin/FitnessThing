import { LeftArrowIcon, ElipsisIcon } from "@/config/icons";
import { Dropdown, DropdownTrigger, DropdownMenu } from "@nextui-org/react";
import { useRouter } from "next/navigation";

export function PageHeader({
    children,
    title,
    enableBackButton,
}: {
    children?: React.ReactElement[] | React.ReactElement;
    title: string;
    enableBackButton?: boolean;
}) {
    const router = useRouter();

    return (
        <section className="flex flex-row items-start justify-between gap-2 px-4">
            {enableBackButton && (
                <div className="h-8 items-center justify-center flex">
                    <LeftArrowIcon
                        className="w-5 h-5 cursor-pointer"
                        onClick={() => router.back()}
                        cursor="pointer"
                    />
                </div>
            )}
            <h1 className="text-2xl font-bold w-full">{title}</h1>
            <div className="h-8 items-center justify-center flex">
                {children && (
                    <Dropdown>
                        <DropdownTrigger>
                            <ElipsisIcon className="w-6 h-6" cursor="pointer" />
                        </DropdownTrigger>
                        <DropdownMenu aria-label="menu">
                            {children}
                        </DropdownMenu>
                    </Dropdown>
                )}
            </div>
        </section>
    );
}
