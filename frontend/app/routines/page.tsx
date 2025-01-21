"use client";
import { authApi, errUnauthorized } from "@/api/api";
import { Loading } from "@/components/loading";
import { PageHeader } from "@/components/page-header";
import { LeftArrowIcon, PlusIcon } from "@/config/icons";
import { Button } from "@nextui-org/button";
import { Card, CardBody, CardFooter, CardHeader } from "@nextui-org/card";
import {
    Form,
    Input,
    Link,
    Modal,
    ModalContent,
    ModalHeader,
    useDisclosure,
} from "@nextui-org/react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function RoutinesPage() {
    const [isLoading, setIsLoading] = useState(true);
    const [isError, setIsError] = useState(false);

    const { isOpen, onOpen, onOpenChange } = useDisclosure();

    const [routines, setRoutines] = useState<any[]>([]);

    const router = useRouter();

    async function fetchData() {
        setIsLoading(true);
        authApi.v1
            .routineServiceGetRoutines()
            .then((response) => {
                console.log(response.data);
                setIsError(false);
                setRoutines(response.data.routines!);
            })
            .catch((error) => {
                console.log(error);
                if (error === errUnauthorized) {
                    router.push("/auth/login");
                    return;
                }
                setIsError(true);
            })
            .finally(() => {
                setIsLoading(false);
            });
    }

    useEffect(() => {
        fetchData();
    }, []);

    if (isLoading) {
        return <Loading />;
    }

    if (isError) {
        return (
            <div className="p-4">
                <h2 className="text-lg text-red-500">
                    Ошибка при загрузке данных
                </h2>
                <p>Проверьте соединение с сервером или обновите страницу.</p>
            </div>
        );
    }

    return (
        <>
            <div className="py-4 flex-grow">
                <PageHeader title="Шаблоны" enableBackButton />
                <div className="grid grid-cols-1 gap-4 p-4">
                    {routines.map((routine) => (
                        <Card
                            key={routine.id}
                            className="flex flex-row flex-grow p-2 gap-4 justify-between"
                            shadow="sm"
                            fullWidth
                            as={Link}
                            href={`/routines/${routine.id}`}
                        >
                            <div className="flex flex-col items-start justify-between p-2">
                                <CardHeader className="p-0">
                                    <label className="text-m font-bold">
                                        {routine.name}
                                    </label>
                                </CardHeader>
                                <CardBody className="p-0">
                                    {routine.description ? (
                                        <div className="py-2">
                                            <p className="text-xs text-gray-400/80">
                                                {routine.description}
                                            </p>
                                        </div>
                                    ) : null}
                                    <p className="text-xs text-gray-500">
                                        {Math.ceil(
                                            Math.random() * (10 - 5) + 5
                                        )}{" "}
                                        упражнений
                                    </p>
                                </CardBody>
                            </div>
                        </Card>
                    ))}
                    <Button color="primary" onPress={onOpen}>
                        <PlusIcon className="w-4 h-4" />
                        Добавить рутину
                    </Button>
                </div>
            </div>
            <Modal
                isOpen={isOpen}
                onClose={onOpenChange}
                size="xs"
                placement="center"
            >
                <ModalContent>
                    {(onClose) => (
                        <>
                            <ModalHeader>
                                <h2>Добавить рутину</h2>
                            </ModalHeader>
                            <Form
                                validationBehavior="native"
                                className="inline-block text-center justify-center w-full max-w-lg p-4"
                                onSubmit={async (
                                    e: React.FormEvent<HTMLFormElement>
                                ) => {
                                    e.preventDefault();
                                    const data = Object.fromEntries(
                                        new FormData(e.currentTarget)
                                    );

                                    console.log(data);
                                    await authApi.v1
                                        .routineServiceCreateRoutine({
                                            name: data.name.toString(),
                                            description:
                                                data.description.toString(),
                                        })
                                        .then((response) => {
                                            console.log(response.data);
                                            router.push(
                                                `/routines/${response.data.routine?.id}`
                                            );
                                        })
                                        .catch((error) => {
                                            console.log(error);
                                            return error;
                                        });
                                }}
                            >
                                <div className="flex flex-col items-center justify-center gap-4 py-4">
                                    <Input
                                        isRequired
                                        label="Название"
                                        placeholder="Название"
                                        labelPlacement="outside"
                                        name="name"
                                        type="text"
                                        autoFocus
                                    />
                                    <Input
                                        label="Описание"
                                        placeholder="Описание"
                                        labelPlacement="outside"
                                        name="description"
                                        type="text"
                                    />
                                    <Button color="primary" type="submit">
                                        Добавить
                                    </Button>
                                </div>
                            </Form>
                        </>
                    )}
                </ModalContent>
            </Modal>
        </>
    );
}
