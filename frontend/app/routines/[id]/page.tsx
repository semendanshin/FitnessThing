"use client";

import { authApi, errUnauthorized } from "@/api/api";
import { WorkoutRoutineDetailResponse } from "@/api/api.generated";
import { Loading } from "@/components/loading";
import { PageHeader } from "@/components/page-header";
import { ModalSelectExercise } from "@/components/pick-exercises-modal";
import {
    ElipsisIcon,
    LeftArrowIcon,
    PlusIcon,
    TrashCanIcon,
} from "@/config/icons";
import { Button } from "@nextui-org/button";
import { Card, CardBody, CardHeader } from "@nextui-org/card";
import {
    Dropdown,
    DropdownItem,
    DropdownMenu,
    DropdownTrigger,
} from "@nextui-org/dropdown";
import { Form } from "@nextui-org/form";
import { Input, Textarea } from "@nextui-org/input";
import {
    Modal,
    ModalBody,
    ModalContent,
    ModalHeader,
    useDisclosure,
} from "@nextui-org/modal";
import { useRouter } from "next/navigation";
import { use, useEffect, useState } from "react";
import { toast } from "react-toastify";

export default function RoutineDetailsPage({
    params,
}: {
    params: Promise<{ id: string }>;
}) {
    const [isLoading, setIsLoading] = useState(true);

    const [routineDetails, setRoutineDetails] =
        useState<WorkoutRoutineDetailResponse>({});

    const { id } = use(params);

    const { isOpen, onOpen, onClose } = useDisclosure();

    const {
        isOpen: isRenameOpen,
        onOpen: onRenameOpen,
        onOpenChange: onRenameOpenChange,
    } = useDisclosure();

    const router = useRouter();

    async function fetchRoutineDetails() {
        console.log(id);
        await authApi.v1
            .routineServiceGetRoutineDetail(id)
            .then((response) => {
                console.log(response.data);
                setRoutineDetails(response.data);
            })
            .catch((error) => {
                console.log(error);
                if (
                    error === errUnauthorized ||
                    (error as any).response?.status === 401
                ) {
                    router.push("/auth/login");
                    return;
                }
                toast.error("Ошибка при загрузке данных");
            });
    }

    async function addExerciseToRoutine(exerciseId: string) {
        await authApi.v1
            .routineServiceAddExerciseToRoutine(id, { exerciseId: exerciseId })
            .then((response) => {
                console.log(response.data);
            })
            .catch((error) => {
                console.log(error);
                if (
                    error === errUnauthorized ||
                    error.response?.status === 401
                ) {
                    router.push("/auth/login");
                    return;
                }
                toast.error("Ошибка при добавлении упражнения");
            });
    }

    async function submitPickExercise(exerciseIds: string[]) {
        console.log(exerciseIds);
        try {
            setIsLoading(true);
            await Promise.all(
                exerciseIds.map((eId) => addExerciseToRoutine(eId))
            );
            await fetchRoutineDetails();
        } catch (error) {
            console.log(error);
            toast.error("Ошибка при добавлении упражнения");
            if (
                error === errUnauthorized ||
                (error as any).response?.status === 401
            ) {
                router.push("/auth/login");
                return;
            }
        } finally {
            setIsLoading(false);
        }
    }

    async function fetchData() {
        setIsLoading(true);
        await fetchRoutineDetails();
        setIsLoading(false);
    }

    useEffect(() => {
        fetchData();
    }, []);

    if (isLoading) {
        return <Loading />;
    }

    function ExercieInstanceCard({ exerciseInstance }: any) {
        const [isButtonLoading, setIsButtonLoading] = useState(false);

        async function onExerciseInstanceDelete() {
            try {
                setIsButtonLoading(true);
                await authApi.v1.routineServiceRemoveExerciseInstanceFromRoutine(
                    id,
                    exerciseInstance.exerciseInstance.id
                );
                await fetchRoutineDetails();
            } catch (error) {
                console.log(error);
                toast.error("Ошибка при удалении упражнения");
                if (
                    error === errUnauthorized ||
                    (error as any).response?.status === 401
                ) {
                    router.push("/auth/login");
                    return;
                }
            } finally {
                setIsButtonLoading(false);
            }
        }

        return (
            <Card
                key={exerciseInstance.exerciseInstance.id}
                className="flex flex-row flex-grow p-2 gap-4 justify-between"
                shadow="sm"
                fullWidth
            >
                <div className="flex flex-col items-start justify-between p-2">
                    <CardHeader className="p-0">
                        <label className="text-m font-bold">
                            {exerciseInstance.exercise.name}
                        </label>
                    </CardHeader>
                    <CardBody className="p-0">
                        <div className="flex flex-row gap-1">
                            <p className="text-xs text-gray-400 whitespace-nowrap">
                                {exerciseInstance.exerciseInstance.sets
                                    ?.length || 0}{" "}
                                {"подход" +
                                    ((exerciseInstance.exerciseInstance.sets
                                        ?.length || 0) %
                                        10 ===
                                    1
                                        ? ""
                                        : (exerciseInstance.exerciseInstance
                                                .sets?.length || 0) %
                                                10 >=
                                                2 &&
                                            (exerciseInstance.exerciseInstance
                                                .sets?.length || 0) %
                                                10 <=
                                                4
                                          ? "а"
                                          : "ов")}
                                {" •"}
                            </p>
                            <p className="text-xs text-gray-400 whitespace-nowrap">
                                {exerciseInstance.exercise.targetMuscleGroups.join(
                                    ", "
                                )}
                            </p>
                        </div>
                    </CardBody>
                </div>
                <div className="flex flex-col items-center justify-center p-2">
                    <Button
                        color="danger"
                        size="sm"
                        onPress={onExerciseInstanceDelete}
                        isLoading={isButtonLoading}
                        isIconOnly
                    >
                        {isButtonLoading ? (
                            ""
                        ) : (
                            <TrashCanIcon className="w-3 h-3" />
                        )}
                    </Button>
                </div>
            </Card>
        );
    }

    function ModalRenameRoutine() {
        const [routineName, setRoutineName] = useState(
            routineDetails.routine?.name
        );
        const [routineDescription, setRoutineDescription] = useState(
            routineDetails.routine?.description
        );

        async function updateRoutine() {
            // await authApi.v1
            //     .fitnessServiceUpdateRoutine(id, {
            //         name: routineName,
            //         description: routineDescription,
            //     })
            //     .then((response) => {
            //         console.log(response.data);
            //         setIsError(false);
            //         setRoutineDetails(response.data);
            //     })
            //     .catch((error) => {
            //         console.log(error);
            //         if (
            //             error === errUnauthorized ||
            //             error.response?.status === 401
            //         ) {
            //             router.push("/auth/login");
            //             return;
            //         }
            //         setIsError(true);
            //     });
        }

        return (
            <Modal
                isOpen={isRenameOpen}
                onClose={onRenameOpenChange}
                size="xs"
                placement="center"
                scrollBehavior="inside"
                className="overflow-y-auto p-2 w-full"
            >
                <ModalContent>
                    {(onClose) => (
                        <>
                            <ModalHeader className="p-2">
                                Переименовать рутину
                            </ModalHeader>

                            <Form
                                className="inline-block text-center justify-center w-full max-w-lg p-0"
                                onSubmit={(e) => {
                                    e.preventDefault();
                                    try {
                                        updateRoutine();
                                        fetchRoutineDetails();
                                        onClose();
                                    } catch (error) {
                                        console.log(error);
                                        if (
                                            error === errUnauthorized ||
                                            (error as any).response?.status ===
                                                401
                                        ) {
                                            router.push("/auth/login");
                                            return;
                                        }
                                        toast.error(
                                            "Ошибка при обновлении рутины"
                                        );
                                    }
                                }}
                            >
                                <div className="grid grid-cols-1 gap-4 p-2">
                                    <Input
                                        placeholder="Название"
                                        value={routineName}
                                        onChange={(e) =>
                                            setRoutineName(e.target.value)
                                        }
                                        autoFocus
                                    />
                                    <Textarea
                                        placeholder="Описание"
                                        value={
                                            routineDescription
                                                ? routineDescription
                                                : ""
                                        }
                                        onChange={(e) =>
                                            setRoutineDescription(
                                                e.target.value
                                            )
                                        }
                                    />
                                    <Button color="primary" type="submit">
                                        Сохранить
                                    </Button>
                                </div>
                            </Form>
                        </>
                    )}
                </ModalContent>
            </Modal>
        );
    }

    async function onDelete() {
        setIsLoading(true);
        await authApi.v1
            .routineServiceDeleteRoutine(id)
            .catch((error) => {
                console.log(error);
                if (
                    error === errUnauthorized ||
                    (error as any).response?.status === 401
                ) {
                    router.push("/auth/login");
                    return;
                }
                toast.error("Ошибка при удалении рутины");
            })
            .then(() => {
                router.push("/routines");
            });
    }

    return (
        <>
            <div className="flex-grow py-4">
                <PageHeader
                    title={routineDetails.routine?.name!}
                    enableBackButton
                >
                    <DropdownItem
                        key="rename"
                        color="primary"
                        onPress={onRenameOpen}
                    >
                        Переименовать
                    </DropdownItem>
                    <DropdownItem
                        key="delete"
                        className="text-danger"
                        color="danger"
                        onPress={onDelete}
                    >
                        Удалить
                    </DropdownItem>
                </PageHeader>
                {routineDetails.routine?.description ? (
                    <div className="py-2 px-4">
                        <p className="text-sm text-gray-500">
                            {routineDetails.routine?.description}
                        </p>
                    </div>
                ) : null}
                <div className="grid grid-cols-1 gap-4 p-4">
                    {routineDetails.exerciseInstances?.map(
                        (exerciseInstance: any) => (
                            <ExercieInstanceCard
                                key={exerciseInstance.exerciseInstance.id}
                                exerciseInstance={exerciseInstance}
                            />
                        )
                    )}
                    <Button color="primary" onPress={onOpen}>
                        <PlusIcon className="w-4 h-4" />
                        Добавить упражнение
                    </Button>
                </div>
            </div>
            <ModalSelectExercise
                isOpen={isOpen}
                onClose={onClose}
                excludeExerciseIds={routineDetails.exerciseInstances!.map(
                    (ei) => ei.exercise!.id!
                )}
                onSubmit={submitPickExercise}
            />

            <ModalRenameRoutine />
        </>
    );
}
